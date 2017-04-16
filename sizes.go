package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// For now, assume that the files are in a directory named 'zephyr'
// and remove everything up until that point.  A better solution would
// be to either use some known string, or find a common prefix to most
// of the strings.
var reZephyrdir = regexp.MustCompile(`^.*zephyr/`)

func main() {
	groups := make([]Entity, 0)

	for _, fname := range os.Args[1:] {
		fmt.Printf("reading: %q\n", fname)
		sizes, err := getSizes(fname)
		if err != nil {
			panic(err)
		}

		groups = append(groups, Entity{
			Name:  fname,
			Sizes: sizes,
		})
		// report(sizes)
	}

	switch len(groups) {
	case 1:
		report(&groups[0])
	case 2:
		deltaReport(&groups[0], &groups[1])
	default:
		fmt.Printf("Warning: No reports for 3 or more elf files given\n")
	}
}

type Entity struct {
	Name  string
	Sizes []*Symbol
}

func getSizes(name string) ([]*Symbol, error) {
	cmd := exec.Command("arm-none-eabi-nm",
		"-S", "-l", "--size-sort", "--radix=d", name)
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	result := make([]*Symbol, 0)

	buf := bytes.NewBuffer(out)
	for {
		line, err := buf.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		line = line[:len(line)-1]

		// Fields are address, size, kind, and information,
		// where information may have a tab giving extra
		// symbol information.
		fields := strings.SplitN(line, " ", 4)

		var symbol Symbol

		symbol.Info = ParseInfo(fields[3])

		symbol.Address, err = strconv.ParseUint(fields[0], 10, 64)
		if err != nil {
			panic(err)
		}
		symbol.Size, err = strconv.ParseUint(fields[1], 10, 64)
		if err != nil {
			panic(err)
		}
		symbol.Kind = fields[2]

		result = append(result, &symbol)
	}

	sort.Sort(BySymbol(result))

	// fmt.Printf("Read %d symbols\n", len(result))

	return result, nil
}

func report(ent *Entity) {
	sizes := ent.Sizes
	lastFile := "INVALID"

	for _, sym := range sizes {
		if sym.Kind != "r" && sym.Kind != "R" &&
			sym.Kind != "t" && sym.Kind != "T" {
			continue
		}
		if sym.Info.File != lastFile {
			fmt.Printf("%s:\n", sym.Info.File)
			lastFile = sym.Info.File
		}

		fmt.Printf("%6d %s\n", sym.Size, sym.Info.Symbol)
	}
}

type LastState struct {
	started bool
	file    string
}

func (l *LastState) Show(file string) {
	if !l.started || file != l.file {
		fmt.Printf("%q:\n", file)
		l.file = file
		l.started = true
	}
}

type Symbol struct {
	Address uint64
	Size    uint64
	Kind    string
	Info    Info
}

type Info struct {
	Symbol string
	File   string
	Line   int64
}

// Decode the string into symbol information.
func ParseInfo(text string) (result Info) {
	fields := strings.Split(text, "\t")
	if len(fields) == 1 {
		result.Symbol = fields[0]
	} else if len(fields) == 2 {
		result.Symbol = fields[0]

		fl := strings.Split(fields[1], ":")
		if len(fl) != 2 {
			panic("File:line does not have a single tab")
		}
		result.File = reZephyrdir.ReplaceAllLiteralString(fl[0], "")
		var err error
		result.Line, err = strconv.ParseInt(fl[1], 10, 64)
		if err != nil {
			panic(err)
		}
	}

	return
}

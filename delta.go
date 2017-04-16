package main

import (
	"fmt"
	"sort"
)

// Show the delta between two runs.
func deltaReport(a *Entity, b *Entity) {
	aa := a.Sizes
	bb := b.Sizes

	ai := 0
	bi := 0

	changes := make([]*Change, 0)

	anew := func() {
		changes = append(changes, &Change{
			kind:   'D',
			osize:  int64(aa[ai].Size),
			nsize:  -1,
			file:   aa[ai].Info.File,
			symbol: aa[ai].Info.Symbol,
		})
		// fmt.Printf("D %6d        %s\n", aa[ai].Size, aa[ai].Info.Symbol)
		ai++
	}

	bnew := func() {
		changes = append(changes, &Change{
			kind:   'A',
			osize:  -1,
			nsize:  int64(bb[bi].Size),
			file:   bb[bi].Info.File,
			symbol: bb[bi].Info.Symbol,
		})
		// fmt.Printf("A %6d        %s\n", bb[bi].Size, bb[bi].Info.Symbol)
		bi++
	}

	// TODO: We should sort by symbol for comparison, but group by
	// file for printing.
	for ai < len(aa) || bi < len(bb) {
		if bi == len(bb) {
			anew()
			continue
		}

		if ai == len(aa) {
			bnew()
			continue
		}

		// Both are present, determine the ordering of the
		// symbol.
		if aa[ai].Info.Symbol < bb[bi].Info.Symbol {
			anew()
			continue
		}

		if aa[ai].Info.Symbol > bb[bi].Info.Symbol {
			bnew()
			continue
		}

		// Otherwise, the name is the same, print if the size
		// changed.
		if aa[ai].Size != bb[bi].Size {
			changes = append(changes, &Change{
				kind:   '-',
				osize:  int64(aa[ai].Size),
				nsize:  int64(bb[bi].Size),
				file:   bb[bi].Info.File,
				symbol: bb[bi].Info.Symbol,
			})
		}
		ai++
		bi++
	}

	sort.Sort(ByFileName(changes))

	showChanges(changes)
}

type Change struct {
	kind   byte
	osize  int64
	nsize  int64
	file   string
	symbol string
}

func showChanges(changes []*Change) {
	var state LastState

	ototal := int64(0)
	ntotal := int64(0)

	for _, ch := range changes {
		state.Show(ch.file)
		var otext, ntext string
		if ch.osize == -1 {
			otext = "      "
		} else {
			otext = fmt.Sprintf("%6d", ch.osize)
			ototal += ch.osize
		}
		if ch.nsize == -1 {
			ntext = "      "
		} else {
			ntext = fmt.Sprintf("%6d", ch.nsize)
			ntotal += ch.nsize
		}
		fmt.Printf("%c %s %s %s\n", ch.kind, otext, ntext, ch.symbol)
	}

	fmt.Printf("  %6d %6d TOTAL\n", ototal, ntotal)
}

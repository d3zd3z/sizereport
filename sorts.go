package main

// Sort the symbols grouping the files together, and then grouping by
// name within the files.
type ByFileName []*Change

func (p ByFileName) Len() int { return len(p) }
func (p ByFileName) Less(i, j int) bool {
	if p[i].file == p[j].file {
		return p[i].symbol < p[j].symbol
	} else {
		return p[i].file < p[j].file
	}
}
func (p ByFileName) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

type BySymbol []*Symbol

func (p BySymbol) Len() int           { return len(p) }
func (p BySymbol) Less(i, j int) bool { return p[i].Info.Symbol < p[j].Info.Symbol }
func (p BySymbol) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

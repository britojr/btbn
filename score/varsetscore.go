package score

import "github.com/willf/bitset"

// varsetScore a pair (varset, score)
type varsetScore struct {
	scor float64
	vars *bitset.BitSet
}

// varsetScores defines a sortable list of pairs (varset, score)
type varsetScores []varsetScore

func (vs varsetScores) Len() int {
	return len(vs)
}
func (vs varsetScores) Swap(i, j int) {
	vs[i], vs[j] = vs[j], vs[i]
}
func (vs varsetScores) Less(i, j int) bool {
	// we want the max score on first position
	return vs[i].scor > vs[j].scor
}

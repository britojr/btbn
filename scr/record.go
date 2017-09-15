package scr

import "sort"

// Record defines a pair (score, data), where score can be used as a key to sort a slice
type Record struct {
	score float64
	data  interface{}
}

// SortRecord sorts a score record by descending score
func SortRecord(rs []Record) {
	sort.Slice(rs, func(i int, j int) bool {
		return rs[i].score > rs[j].score
	})
}

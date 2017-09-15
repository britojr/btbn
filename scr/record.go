package scr

import "sort"

// Record defines a pair (score, data), where score can be used as a key to sort a slice
type Record struct {
	score float64
	data  interface{}
}

// SortRecords sorts a score record by descending score
func SortRecords(rs []*Record) {
	sort.Slice(rs, func(i int, j int) bool {
		return rs[i].score > rs[j].score
	})
}

// NewRecord returns record data
func NewRecord(score float64, data interface{}) *Record {
	return &Record{score, data}
}

// Data returns record data
func (r *Record) Data() interface{} {
	return r.data
}

// Score returns record score
func (r *Record) Score() float64 {
	return r.score
}

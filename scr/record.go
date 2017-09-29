package scr

import "sort"

// Record defines a pair (score, data) allowing a slice of records to be sorted by score
type Record struct {
	score float64
	data  interface{}
}

// SortRecords sorts a score record by descending score
func SortRecords(rs []*Record, descending bool) {
	f := func(i, j int) bool { return rs[i].score < rs[j].score }
	if descending {
		f = func(i, j int) bool { return rs[i].score > rs[j].score }
	}
	sort.Slice(rs, f)
}

// NewRecord returns pointer to a new record data
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

// RecordHeap defines a heap implementation for a slice of records
type RecordHeap struct {
	recs *[]*Record
	less func(i, j int) bool
}

// NewRecordHeap creates a new record heap
// the less function defines if it is a min or max heap
func NewRecordHeap(recs *[]*Record, less func(i, j int) bool) *RecordHeap {
	return &RecordHeap{recs, less}
}

func (rh RecordHeap) Len() int           { return len(*(rh.recs)) }
func (rh RecordHeap) Swap(i, j int)      { (*rh.recs)[i], (*rh.recs)[j] = (*rh.recs)[j], (*rh.recs)[i] }
func (rh RecordHeap) Less(i, j int) bool { return rh.less(i, j) }

// Push appends a record to record heap
func (rh *RecordHeap) Push(x interface{}) { *rh.recs = append(*rh.recs, x.(*Record)) }

// Pop removes the top record from record heap
func (rh *RecordHeap) Pop() interface{} {
	n := len(*rh.recs) - 1
	rec := (*rh.recs)[n]
	*rh.recs = (*rh.recs)[:n]
	return rec
}

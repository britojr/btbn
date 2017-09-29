package scr

import (
	"container/heap"
	"reflect"
	"testing"
)

func TestRecordHeap(t *testing.T) {
	rs1 := []*Record{{5, nil}, {2, nil}, {7, nil}}
	rs2 := []*Record{{5, nil}, {2, nil}, {7, nil}}
	cases := []struct {
		init    *[]*Record
		less    func(i, j int) bool
		actions []func(*RecordHeap)
		want    []float64
	}{{
		init: &rs1,
		less: func(i, j int) bool { return rs1[i].score < rs1[j].score },
		actions: []func(*RecordHeap){
			func(h *RecordHeap) { heap.Push(h, NewRecord(9, nil)) },
			func(h *RecordHeap) { heap.Push(h, NewRecord(3, nil)) },
			func(h *RecordHeap) { heap.Pop(h) },
			func(h *RecordHeap) { heap.Push(h, NewRecord(1, nil)) },
			func(h *RecordHeap) { heap.Pop(h) },
		},
		want: []float64{3, 5, 7, 9},
	}, {
		init: &rs2,
		less: func(i, j int) bool { return rs2[i].score > rs2[j].score },
		actions: []func(*RecordHeap){
			func(h *RecordHeap) { heap.Push(h, NewRecord(9, nil)) },
			func(h *RecordHeap) { heap.Push(h, NewRecord(3, nil)) },
			func(h *RecordHeap) { heap.Pop(h) },
			func(h *RecordHeap) { heap.Push(h, NewRecord(1, nil)) },
			func(h *RecordHeap) { heap.Pop(h) },
		},
		want: []float64{5, 3, 2, 1},
	}}

	for _, tt := range cases {
		h := NewRecordHeap(tt.init, tt.less)
		heap.Init(h)
		for _, run := range tt.actions {
			run(h)
		}
		got := helperGetPopKeys(h)
		if !reflect.DeepEqual(tt.want, got) {
			t.Errorf("heap not working (%v) != (%v)", tt.want, got)
		}
	}
}

func helperGetPopKeys(h *RecordHeap) (ks []float64) {
	for h.Len() > 0 {
		ks = append(ks, heap.Pop(h).(*Record).Score())
	}
	return
}

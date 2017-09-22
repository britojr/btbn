package optimizer

import (
	"fmt"
	"reflect"
	"sort"
	"testing"
)

func TestGetInitialDAG(t *testing.T) {
	cases := []struct {
		iniclq []int
		remain []int
	}{{
		[]int{0, 1, 2},
		[]int{3, 4, 5, 6, 7, 8, 9},
	}, {
		[]int{4, 7, 3},
		[]int{0, 1, 2, 5, 6, 8, 9},
	}}
	for _, tt := range cases {
		n := len(tt.iniclq) + len(tt.remain)
		s := &IterativeSearch{common: newCommon(&fakeRanker{n})}
		got := s.getInitialDAG(tt.iniclq)
		if n != got.Size() {
			t.Errorf("wrong bn size, want (%v), got (%v)", n, got)
		}
		for _, v := range tt.remain {
			if got.Parents(v).Count() != 0 {
				t.Errorf("set parent for variable (%v) not in initial clique (%v)", v, tt.iniclq)
			}
		}
		sum := 0
		for _, v := range tt.iniclq {
			sum += got.Parents(v).Count()
		}
		if sum == 0 {
			t.Errorf("initial clique (%v) is empty", tt.iniclq)
		}
	}
}

func TestSampleOrder(t *testing.T) {
	cases := []struct {
		n, k int
		prev map[string]struct{}
		only []int
	}{{
		5, 2, map[string]struct{}{
			fmt.Sprint([]int{0, 1, 2}): struct{}{},
			fmt.Sprint([]int{0, 1, 3}): struct{}{},
			fmt.Sprint([]int{0, 1, 4}): struct{}{},
			fmt.Sprint([]int{0, 2, 3}): struct{}{},
			fmt.Sprint([]int{0, 2, 4}): struct{}{},
			fmt.Sprint([]int{0, 3, 4}): struct{}{},
			fmt.Sprint([]int{1, 2, 3}): struct{}{},
			fmt.Sprint([]int{1, 3, 4}): struct{}{},
			fmt.Sprint([]int{2, 3, 4}): struct{}{},
		},
		[]int{1, 2, 4},
	}}
	for _, tt := range cases {
		s := &IterativeSearch{common: newCommon(&fakeRanker{tt.n})}
		s.prevCliques = tt.prev
		s.tw = tt.k
		got := s.sampleOrder()
		if tt.n != len(got) {
			t.Errorf("wrong order size, want (%v), got (%v)", tt.n, got)
		}
		if len(tt.prev) > 0 {
			ini := got[:tt.k+1]
			sort.Ints(ini)
			if !reflect.DeepEqual(tt.only, ini) {
				t.Errorf("didn't sample the only possible (k+1)-clique (%v), got (%v)", tt.only, ini)
			}
		}
	}
}

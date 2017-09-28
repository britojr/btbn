package optimizer

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/britojr/btbn/scr"
	"github.com/britojr/btbn/varset"
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
			fmt.Sprint([]int{3, 4}): struct{}{},
			fmt.Sprint([]int{4, 3}): struct{}{},
			fmt.Sprint([]int{2, 4}): struct{}{},
			fmt.Sprint([]int{4, 2}): struct{}{},
			fmt.Sprint([]int{1, 4}): struct{}{},
			fmt.Sprint([]int{4, 1}): struct{}{},
			fmt.Sprint([]int{0, 4}): struct{}{},
			fmt.Sprint([]int{4, 0}): struct{}{},
			fmt.Sprint([]int{2, 3}): struct{}{},
			fmt.Sprint([]int{3, 2}): struct{}{},
			fmt.Sprint([]int{1, 3}): struct{}{},
			fmt.Sprint([]int{3, 1}): struct{}{},
			fmt.Sprint([]int{0, 3}): struct{}{},
			// fmt.Sprint([]int{3,0}): struct{}{},
			fmt.Sprint([]int{1, 2}): struct{}{},
			fmt.Sprint([]int{2, 1}): struct{}{},
			fmt.Sprint([]int{0, 2}): struct{}{},
			fmt.Sprint([]int{2, 0}): struct{}{},
			fmt.Sprint([]int{0, 1}): struct{}{},
			fmt.Sprint([]int{1, 0}): struct{}{},
		},
		[]int{3, 0},
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
			ini := got[tt.k+1:]
			if !reflect.DeepEqual(tt.only, ini) {
				t.Errorf("didn't sample the only possible order (%v), got (%v)", tt.only, ini)
			}
		}
	}
}

func TestGreedySearch(t *testing.T) {
	fname := helperGetTempFile(fixedScores, "iter_test")
	defer os.Remove(fname)
	rk := scr.CreateRanker(scr.Read(fname), 10)
	cases := []struct {
		ranker  scr.Ranker
		initps  map[int]varset.Varset
		ord     []int
		n, k    int
		wantps  map[int]varset.Varset
		wantscr float64
	}{{
		ranker: rk, ord: []int{0, 1, 2, 3, 4, 5, 6}, n: rk.Size(), k: 2,
		initps: map[int]varset.Varset{
			0: varset.New(rk.Size()),
			1: varset.New(rk.Size()).SetInts([]int{0, 2}),
			2: varset.New(rk.Size()),
		},
		wantps: map[int]varset.Varset{
			0: varset.New(rk.Size()),
			1: varset.New(rk.Size()).SetInts([]int{0, 2}),
			2: varset.New(rk.Size()),
			3: varset.New(rk.Size()).SetInts([]int{0, 2}),
			4: varset.New(rk.Size()).SetInts([]int{0, 2}),
			5: varset.New(rk.Size()).SetInts([]int{0, 1}),
			6: varset.New(rk.Size()).SetInts([]int{5, 1}),
		},
		wantscr: -380,
	}}
	for _, tt := range cases {
		bn := helperGetBN(tt.initps, tt.ranker)
		s := &IterativeSearch{common: newCommon(tt.ranker)}
		s.tw, s.nv = tt.k, tt.n
		s.greedySearch(bn, tt.ord)
		for i := 0; i < tt.ranker.Size(); i++ {
			if !tt.wantps[i].Equal(bn.Parents(i)) {
				t.Errorf("wrong parents for %v: (%v)!=(%v)", i, tt.wantps[i], bn.Parents(i))
			}
		}
		if tt.wantscr != bn.Score() {
			t.Errorf("wrong total score (%v)!=(%v)", tt.wantscr, bn.Score())
		}
	}
}

func TestAstarSearch(t *testing.T) {
	fname := helperGetTempFile(fixedScores, "iter_test")
	defer os.Remove(fname)
	rk := scr.CreateRanker(scr.Read(fname), 10)
	cases := []struct {
		ranker  scr.Ranker
		initps  map[int]varset.Varset
		ord     []int
		n, k    int
		wantps  map[int]varset.Varset
		wantscr float64
	}{{
		ranker: rk, ord: []int{0, 1, 2, 3, 4, 5, 6}, n: rk.Size(), k: 2,
		initps: map[int]varset.Varset{
			0: varset.New(rk.Size()),
			1: varset.New(rk.Size()).SetInts([]int{0, 2}),
			2: varset.New(rk.Size()),
		},
		wantps: map[int]varset.Varset{
			0: varset.New(rk.Size()),
			1: varset.New(rk.Size()).SetInts([]int{0, 2}),
			2: varset.New(rk.Size()),
			3: varset.New(rk.Size()).SetInts([]int{0, 2}),
			4: varset.New(rk.Size()).SetInts([]int{3, 2}),
			5: varset.New(rk.Size()).SetInts([]int{0, 1}),
			6: varset.New(rk.Size()).SetInts([]int{3, 4}),
		},
		wantscr: -350,
	}}
	for _, tt := range cases {
		bn := helperGetBN(tt.initps, tt.ranker)
		s := &IterativeSearch{common: newCommon(tt.ranker)}
		s.tw, s.nv = tt.k, tt.n
		s.astarSearch(bn, tt.ord)
		for i := 0; i < tt.ranker.Size(); i++ {
			if !tt.wantps[i].Equal(bn.Parents(i)) {
				t.Errorf("wrong parents for %v: (%v)!=(%v)", i, tt.wantps[i], bn.Parents(i))
			}
		}
		if tt.wantscr != bn.Score() {
			t.Errorf("wrong total score (%v)!=(%v)", tt.wantscr, bn.Score())
		}
	}
}

func helperGetBN(psets map[int]varset.Varset, ranker scr.Ranker) *BNStructure {
	bn := NewBNStructure(ranker.Size())
	emp := varset.New(ranker.Size())
	for i := 0; i < bn.Size(); i++ {
		if pset, ok := psets[i]; ok {
			bn.SetParents(i, pset, ranker.ScoreOf(i, pset))
		} else {
			bn.SetParents(i, emp, ranker.ScoreOf(i, emp))
		}
	}
	return bn
}

var fixedScores = `META pss_version = 0.1
VAR 0
-100

VAR 1
-100
-20 0 2

VAR 2
-100
-10 0

VAR 3
-100
-80 0 1
-20 0 2
-80 1 2
-10 0 1 2

VAR 4
-100
-99 2
-80 0 1
-60 0 2
-80 1 2
-71 3 0
-70 3 1
-70 3 2
-11 0 1 2
-10 0 1 2 3

VAR 5
-100
-20 0 1
-80 0 2
-80 1 2
-80 3 0
-80 3 1
-80 3 2
-10 0 1 2

VAR 6
-100
-70 0 1
-80 0 2
-80 1 2
-80 3 0
-20 4 3
-80 3 2
-60 5 1
-50 5 2
-10 0 1 2

`

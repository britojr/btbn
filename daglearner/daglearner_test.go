package daglearner

import (
	"reflect"
	"testing"

	"github.com/britojr/btbn/bnstruct"
	"github.com/britojr/btbn/ktree"
	"github.com/britojr/btbn/scr"
	"github.com/britojr/btbn/varset"
)

type fakeRanker struct {
	n int
}

func (r *fakeRanker) BestIn(v int, restric varset.Varset) (varset.Varset, float64) {
	ps := varset.New(0).SetInts(restric.DumpAsInts())
	return ps, r.ScoreOf(v, ps)
}
func (r *fakeRanker) BestInLim(v int, restric varset.Varset, maxPa int) (varset.Varset, float64) {
	ps := varset.New(0).SetInts(restric.DumpAsInts()[:maxPa])
	return ps, r.ScoreOf(v, ps)
}

func (r *fakeRanker) ScoreOf(v int, parents varset.Varset) float64 {
	c := 0
	for _, u := range parents.DumpAsInts() {
		if u > v {
			c++
		}
	}
	return float64(c)
}
func (r *fakeRanker) Size() int {
	return r.n
}

func TestApproxLearning(t *testing.T) {
	cases := []struct {
		n, k   int
		ranker scr.Ranker
	}{
		{7, 2, &fakeRanker{7}},
		{11, 4, &fakeRanker{11}},
	}
	for _, tt := range cases {
		tk := ktree.UniformSample(tt.n, tt.k)
		got := Approximated(tk, tt.ranker)
		clqset := helperTreeVarsets(tk, tt.n)
		for i := 0; i < got.Size(); i++ {
			found := helperIsSubSetOfSome(got.Parents(i).Clone().Set(i), clqset)
			if !found {
				t.Errorf("Clique (%v) not found in ktree (%v)",
					got.Parents(i).Clone().Set(i), tk,
				)
			}
		}
	}
}

func TestSamplePartialOrder(t *testing.T) {
	cases := []struct {
		n, k int
	}{
		{7, 2},
		{11, 4},
	}
	for _, tt := range cases {
		tk := ktree.UniformSample(tt.n, tt.k)
		got := samplePartialOrder(tk)
		clqset := helperTreeVarsets(tk, tt.n)
		for i := range clqset {
			found := helperIsEqualToSome(varset.New(tt.n).SetInts(got[i].vars), clqset)
			if !found {
				t.Errorf("Clique (%v) not found in ktree (%v)", got[i].vars, tk)
			}
		}
	}
}

func TestSetParentsFromOrder(t *testing.T) {
	cases := []struct {
		ranker  scr.Ranker
		order   partialOrder
		parents map[int][]int
	}{{
		&fakeRanker{9},
		partialOrder{[]int{1, 2, 3, 4}, 0},
		map[int][]int{
			1: {},
			2: {1},
			3: {1, 2},
			4: {1, 2, 3},
		}}, {
		&fakeRanker{9},
		partialOrder{[]int{5, 2, 1, 7}, 2},
		map[int][]int{
			5: {},
			2: {},
			1: {2, 5},
			7: {1, 2, 5},
		}}}
	for _, tt := range cases {
		bn := bnstruct.New(tt.ranker.Size())
		setParentsFromOrder(tt.order, tt.ranker, bn)
		for v, ps := range tt.parents {
			if bn.Parents(v) != nil {
				if !reflect.DeepEqual(ps, bn.Parents(v).DumpAsInts()) {
					t.Errorf("wrong parent set for %v: (%v)!=(%v)", v, ps, bn.Parents(v))
				}
			} else if len(ps) > 0 {
				t.Errorf("wrong parent set for %v: (%v)!=(%v)", v, ps, bn.Parents(v))
			}
		}
	}
}

func helperTreeVarsets(tk *ktree.Ktree, n int) (vl []varset.Varset) {
	queue := []*ktree.Ktree{tk}
	for len(queue) > 0 {
		r := queue[0]
		queue = queue[1:]
		vl = append(vl, varset.New(n).SetInts(r.Variables()))
		queue = append(queue, r.Children()...)
	}
	return
}
func helperIsSubSetOfSome(pset varset.Varset, clqset []varset.Varset) bool {
	for _, clq := range clqset {
		if clq.IsSuperSet(pset) {
			return true
		}
	}
	return false
}
func helperIsEqualToSome(pset varset.Varset, clqset []varset.Varset) bool {
	for _, clq := range clqset {
		if clq.Equal(pset) {
			return true
		}
	}
	return false
}

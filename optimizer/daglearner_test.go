package optimizer

import (
	"reflect"
	"testing"

	"github.com/britojr/btbn/ktree"
	"github.com/britojr/btbn/scr"
	"github.com/britojr/btbn/varset"
)

type fakeRanker struct{}

func (m *fakeRanker) BestIn(restric varset.Varset) (parents varset.Varset, scr float64) {
	return m.BestInLim(restric, restric.Count())
}
func (m *fakeRanker) BestInLim(restric varset.Varset, maxPa int) (parents varset.Varset, scr float64) {
	ps := varset.New(0).SetInts(restric.DumpAsInts()[:maxPa])
	return ps, float64(ps.Count())
}
func (m *fakeRanker) ScoreOf(parents varset.Varset) float64 {
	return float64(parents.Count())
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
func TestApproxLearning(t *testing.T) {
	fr := new(fakeRanker)
	cases := []struct {
		n, k    int
		rankers []scr.Ranker
	}{
		{7, 2, []scr.Ranker{fr, fr, fr, fr, fr, fr, fr}},
		{11, 4, []scr.Ranker{fr, fr, fr, fr, fr, fr, fr, fr, fr, fr, fr}},
	}
	for _, tt := range cases {
		tk := ktree.UniformSample(tt.n, tt.k)
		got := DAGapproximatedLearning(tk, tt.rankers)
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

func helperIsEqualToSome(pset varset.Varset, clqset []varset.Varset) bool {
	for _, clq := range clqset {
		if clq.Equal(pset) {
			return true
		}
	}
	return false
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
	fr := new(fakeRanker)
	cases := []struct {
		rankers []scr.Ranker
		order   partialOrder
		parents map[int][]int
	}{{
		[]scr.Ranker{fr, fr, fr, fr, fr, fr, fr, fr, fr},
		partialOrder{[]int{1, 2, 3, 4}, 0},
		map[int][]int{
			1: []int{},
			2: []int{1},
			3: []int{1, 2},
			4: []int{1, 2, 3},
		}}, {
		[]scr.Ranker{fr, fr, fr, fr, fr, fr, fr, fr, fr},
		partialOrder{[]int{5, 2, 1, 7}, 2},
		map[int][]int{
			5: []int{},
			2: []int{},
			1: []int{2, 5},
			7: []int{1, 2, 5},
		}}}
	for _, tt := range cases {
		bn := NewBNStructure(len(tt.rankers))
		setParentsFromOrder(tt.order, tt.rankers, bn)
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

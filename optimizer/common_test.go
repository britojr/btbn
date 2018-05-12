package optimizer

import (
	"os"
	"reflect"
	"testing"

	"github.com/britojr/btbn/scr"
	"github.com/britojr/btbn/varset"
	"github.com/britojr/utl/ioutl"
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
func (r *fakeRanker) SaveSubSet(fname string, vs []int) {}
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

func TestCreate(t *testing.T) {
	ranker := &fakeRanker{5}
	parms := make(map[string]string)
	alg := Create(AlgSampleSearch, ranker, parms)
	if reflect.TypeOf(alg) != reflect.TypeOf(&SampleSearch{}) {
		t.Errorf("wrong optimizer type: (SampleSearch)!=(%v)", reflect.TypeOf(alg))
	}
	alg = Create(AlgIterativeSearch, ranker, parms)
	if reflect.TypeOf(alg) != reflect.TypeOf(&IterativeSearch{}) {
		t.Errorf("wrong optimizer type: (SampleSearch)!=(%v)", reflect.TypeOf(alg))
	}
}

func TestCommonSearch(t *testing.T) {
	mi := `0.50
0.05 0.50
0.11 0.11 0.67
0.11 0.11 0.67 0.70
`
	mifile := ioutl.TempFile("selec_test", mi)
	defer os.Remove(mifile)
	cases := []struct {
		alg    string
		ranker scr.Ranker
		parms  map[string]string
	}{{
		AlgSampleSearch, &fakeRanker{5}, map[string]string{
			ParmTreewidth: "1",
		},
	}, {
		AlgIterativeSearch, &fakeRanker{10}, map[string]string{
			ParmTreewidth:       "4",
			ParmSearchVariation: OpGreedy,
			ParmInitIters:       "50",
		},
	}, {
		AlgSelectedSample, &fakeRanker{4}, map[string]string{
			ParmTreewidth:       "2",
			ParmSearchVariation: OpGreedy,
			ParmNumTrees:        "5",
			ParmMutualInfo:      mifile,
		},
	}}
	for _, tt := range cases {
		bn := Create(tt.alg, tt.ranker, tt.parms).Search()
		if tt.ranker.Size() != bn.Size() {
			t.Errorf("bn with wrong number of vars (%v)!=(%v)", tt.ranker.Size(), bn.Size())
		}
	}
	// TODO: add functions to check cyclicity, connectivity and treewidth
}

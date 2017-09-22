package optimizer

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"github.com/britojr/btbn/scr"
	"github.com/britojr/utl/errchk"
)

func helperGetTempFile(content string, fprefix string) string {
	f, err := ioutil.TempFile("", fprefix)
	errchk.Check(err, "")
	defer f.Close()
	fmt.Fprintf(f, "%s", content)
	return f.Name()
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
	mifile := helperGetTempFile(mi, "selec_test")
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

package scr

import (
	"os"
	"reflect"
	"testing"

	"github.com/britojr/btbn/varset"
)

func TestScoreOf(t *testing.T) {
	p := make(map[string]varset.Varset)
	s := []string{
		"000", "001", "010", "011", "100", "101", "110", "111",
	}
	for i := range s {
		p[s[i]] = varset.New(len(s[i]))
		p[s[i]].LoadHashString(s[i])
	}
	content := `META pss_version = 0.1
META teste 0,1,2,3,4
META num_variables=5

VAR A
-2
-9 B
-8 C
-6 B C

VAR B
-9
-9 A
-1 C
-3 A C

VAR C
-10
-10.1 A
-2 B
-1 A B
`
	cases := []struct {
		content  string
		maxPa, x int
		vars     []int
		scor     float64
	}{
		{content, 0, 0, []int{}, -2},
		{content, 0, 0, []int{1}, -9},
		{content, 0, 0, []int{2}, -8},
		{content, 0, 0, []int{1, 2}, -6},

		{content, 0, 2, []int{}, -10},
		{content, 0, 2, []int{0}, -10.1},
		{content, 0, 2, []int{1}, -2},
		{content, 0, 2, []int{0, 1}, -1},
	}

	for _, tt := range cases {
		fname := helperGetTempFile(tt.content, "rk_test")
		defer os.Remove(fname)
		sr := CreateRanker(Read(fname), tt.maxPa)
		parents := varset.New(sr.Size())
		for _, v := range tt.vars {
			parents.Set(v)
		}
		got := sr.ScoreOf(tt.x, parents)
		if tt.scor != got {
			t.Errorf("wrong scores (%v)!=(%v)", tt.scor, got)
		}
	}
}

func TestBestIn(t *testing.T) {
	content := `META pss_version = 0.1
META teste 0,1,2,3,4
META num_variables=5

VAR A
-2
-9 B
-8 C
-6 B C

VAR B
-9
-9 A
-1 C
-3 A C

VAR C
-10
-10.1 A
-2 B
-1 A B
`
	cases := []struct {
		content  string
		maxPa, x int
		restric  []int
		wantPa   []int
		wantSc   float64
	}{
		{content, 0, 0, []int{0, 1, 2}, []int{}, -2},
		{content, 0, 0, []int{}, []int{}, -2},
		{content, 0, 1, []int{0, 1, 2}, []int{2}, -1},
		{content, 0, 2, []int{0, 1, 2}, []int{0, 1}, -1},
		{content, 0, 2, []int{1, 2}, []int{1}, -2},
		{content, 0, 2, []int{0, 2}, []int{}, -10},
	}

	for _, tt := range cases {
		fname := helperGetTempFile(tt.content, "rk_test")
		defer os.Remove(fname)
		sr := CreateRanker(Read(fname), tt.maxPa)
		restric := varset.New(sr.Size())
		for _, v := range tt.restric {
			restric.Set(v)
		}
		gotPa, gotSc := sr.BestIn(tt.x, restric)
		if tt.wantSc != gotSc {
			t.Errorf("wrong scores (%v)!=(%v)", tt.wantSc, gotSc)
		}
		if !reflect.DeepEqual(tt.wantPa, gotPa.DumpAsInts()) {
			t.Errorf("wrong parents (%v)!=(%v) for x=%v", tt.wantPa, gotPa.DumpAsInts(), tt.x)
		}
	}
}

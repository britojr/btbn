package scr

import (
	"os"
	"testing"

	"github.com/britojr/btbn/varset"
)

func TestScores(t *testing.T) {
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
		x        int
		pxScores []struct {
			vars []int
			scor float64
		}
	}{{
		content, 0,
		[]struct {
			vars []int
			scor float64
		}{
			{[]int{}, -2},
			{[]int{1}, -9},
			{[]int{2}, -8},
			{[]int{1, 2}, -6},
		}}, {
		content, 2,
		[]struct {
			vars []int
			scor float64
		}{
			{[]int{}, -10},
			{[]int{0}, -10.1},
			{[]int{1}, -2},
			{[]int{0, 1}, -1},
		}}}

	for _, tt := range cases {
		fname := helperGetTempFile(tt.content, "rk_test")
		defer os.Remove(fname)
		sc := Read(fname)
		scoreMap := sc.Scores(tt.x)
		for _, px := range tt.pxScores {
			parents := varset.New(sc.Nvar())
			for _, v := range px.vars {
				parents.Set(v)
			}
			if px.scor != scoreMap[parents.DumpHashString()] {
				t.Errorf("wrong scores (%v)!=(%v)", px.scor, scoreMap[parents.DumpHashString()])
			}
		}
	}
}

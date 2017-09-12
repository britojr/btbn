package score

import (
	"testing"

	"github.com/britojr/btbn/varset"
)

func TestScores(t *testing.T) {
	cases := []struct {
		scoreFile string
		x         int
		pxScores  []struct {
			vars []int
			scor float64
		}
	}{
		{"scorefile.test", 0,
			[]struct {
				vars []int
				scor float64
			}{
				{[]int{}, -2},
				{[]int{1}, -9},
				{[]int{2}, -8},
				{[]int{1, 2}, -6},
			},
		}, {
			"scorefile.test", 2,
			[]struct {
				vars []int
				scor float64
			}{
				{[]int{}, -10},
				{[]int{0}, -10},
				{[]int{1}, -2},
				{[]int{0, 1}, -1},
			},
		},
	}

	for _, tt := range cases {
		sc := Read(tt.scoreFile)
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

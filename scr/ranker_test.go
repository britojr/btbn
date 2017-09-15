package scr

import (
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

	cases := []struct {
		scoreFile string
		maxPa, x  int
		vars      []int
		scor      float64
	}{
		{"scorefile.test", 0, 0, []int{}, -2},
		{"scorefile.test", 0, 0, []int{1}, -9},
		{"scorefile.test", 0, 0, []int{2}, -8},
		{"scorefile.test", 0, 0, []int{1, 2}, -6},

		{"scorefile.test", 0, 2, []int{}, -10},
		{"scorefile.test", 0, 2, []int{0}, -10},
		{"scorefile.test", 0, 2, []int{1}, -2},
		{"scorefile.test", 0, 2, []int{0, 1}, -1},
	}

	for _, tt := range cases {
		sr := CreateRankers(Read(tt.scoreFile), tt.maxPa)
		parents := varset.New(len(sr))
		for _, v := range tt.vars {
			parents.Set(v)
		}
		got := sr[tt.x].ScoreOf(parents)
		if tt.scor != got {
			t.Errorf("wrong scores (%v)!=(%v)", tt.scor, got)
		}
	}
}

func TestBestIn(t *testing.T) {
	cases := []struct {
		scoreFile string
		maxPa, x  int
		restric   []int
		wantPa    []int
		wantSc    float64
	}{
		{"scorefile.test", 0, 0, []int{0, 1, 2}, []int{}, -2},
		{"scorefile.test", 0, 0, []int{}, []int{}, -2},
		{"scorefile.test", 0, 1, []int{0, 1, 2}, []int{2}, -1},
		{"scorefile.test", 0, 2, []int{0, 1, 2}, []int{0, 1}, -1},
		{"scorefile.test", 0, 2, []int{1, 2}, []int{1}, -2},
		{"scorefile.test", 0, 2, []int{0, 2}, []int{}, -10},
	}

	for _, tt := range cases {
		sr := CreateRankers(Read(tt.scoreFile), tt.maxPa)
		restric := varset.New(len(sr))
		for _, v := range tt.restric {
			restric.Set(v)
		}
		gotPa, gotSc := sr[tt.x].BestIn(restric)
		if tt.wantSc != gotSc {
			t.Errorf("wrong scores (%v)!=(%v)", tt.wantSc, gotSc)
		}
		if !reflect.DeepEqual(tt.wantPa, gotPa.DumpAsInts()) {
			t.Errorf("wrong parents (%v)!=(%v) for x=%v", tt.wantPa, gotPa.DumpAsInts(), tt.x)
		}
	}
}

package score

import (
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
		p[s[i]].SetFromString(s[i])
	}

	cases := []struct {
		scoreFile string
		maxPa, x  int
		pxScores  map[varset.Varset]float64
	}{
		{
			"scorefile.test",
			0, 0,
			map[varset.Varset]float64{
				p["000"]: -2,
				p["010"]: -9,
				p["001"]: -8,
				p["011"]: -6,
			},
		}, {
			"scorefile.test",
			0, 2,
			map[varset.Varset]float64{
				p["000"]: -10,
				p["100"]: -10,
				p["010"]: -2,
				p["110"]: -1,
			},
		},
	}

	for _, tt := range cases {
		sr := CreateRankers(Read(tt.scoreFile), tt.maxPa)
		for px, scor := range tt.pxScores {
			got := sr[tt.x].ScoreOf(px)
			if scor != got {
				t.Errorf("wrong scores (%v)!=(%v)", scor, got)
			}
		}
	}
}

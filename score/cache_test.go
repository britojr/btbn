package score

import (
	"testing"

	"github.com/britojr/btbn/varset"
)

func TestScores(t *testing.T) {
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
		x         int
		pxScores  map[string]float64
	}{
		{
			"scorefile.test",
			0,
			map[string]float64{
				p["000"].DumpAsString(): -2,
				p["010"].DumpAsString(): -9,
				p["001"].DumpAsString(): -8,
				p["011"].DumpAsString(): -6,
			},
		}, {
			"scorefile.test",
			2,
			map[string]float64{
				p["000"].DumpAsString(): -10,
				p["100"].DumpAsString(): -10,
				p["010"].DumpAsString(): -2,
				p["110"].DumpAsString(): -1,
			},
		},
	}

	for _, tt := range cases {
		sc := Read(tt.scoreFile)
		scoreMap := sc.Scores(tt.x)
		for px := range scoreMap {
			if tt.pxScores[px] != scoreMap[px] {
				t.Errorf("wrong scores (%v)!=(%v)", tt.pxScores[px], scoreMap[px])
			}
		}
	}
}

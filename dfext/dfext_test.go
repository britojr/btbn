package dfext

import (
	"reflect"
	"sort"
	"testing"

	"github.com/kniren/gota/dataframe"
)

func TestCount(t *testing.T) {
	cases := []struct {
		df        dataframe.DataFrame
		normalize bool
		sorted    bool
		want      []float64
	}{{
		df: dataframe.LoadStructs([]struct{ A, B int }{
			{0, 0},
			{1, 0},
			{0, 1},
			{1, 1},
		}),
		normalize: true,
		sorted:    false,
		want:      []float64{.25, .25, .25, .25},
	}, {
		df: dataframe.LoadStructs([]struct{ B, C, D int }{
			{0, 0, 0},
			{1, 0, 0},
			{0, 1, 0},
			{0, 1, 0},
			{1, 0, 0},
			{1, 1, 1},
			{1, 0, 0},
			{1, 1, 0},
		}),
		normalize: true,
		sorted:    false,
		want:      []float64{1.0 / 8.0, 3.0 / 8.0, 2.0 / 8.0, 1.0 / 8.0, 1.0 / 8.0},
	}}
	for _, tt := range cases {
		got := Counts(tt.df, tt.normalize)
		if !tt.sorted {
			sort.Float64s(got)
			sort.Float64s(tt.want)
		}
		if !reflect.DeepEqual(tt.want, got) {
			t.Errorf("wrong count %v != %v", tt.want, got)
		}
	}
}

// Package dfext provides some extentions to dataframe package
package dfext

import (
	"strings"

	"github.com/kniren/gota/dataframe"
)

// Counts returns the occurrence counts of unique joints of values
// the order of the counts is random
func Counts(df dataframe.DataFrame, normalize bool) []float64 {
	// TODO: improve this with functionalities like like pandas:
	// Series.value_counts(normalize=False, sort=True, ascending=False, bins=None, dropna=True)
	c := make(map[string]float64)
	for _, v := range df.Records() {
		c[strings.Join(v, ",")]++
	}
	vs := make([]float64, 0, len(c))
	for _, v := range c {
		if normalize {
			v = v / float64(df.Nrow())
		}
		vs = append(vs, v)
	}
	return vs
}

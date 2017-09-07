package optimizer

import (
	"fmt"

	"github.com/britojr/btbn/score"
)

// Optimizer defines a structure optimizer algorithm
type Optimizer interface {
	Search(numSolutions, timeAvailable int) interface{}
}

// Create creates a structure optimizer algorithm
func Create(optimizerAlg string, scoreRankers []score.Ranker, parmFile string) Optimizer {
	switch optimizerAlg {
	case "sample":
		return NewSampleSearch(scoreRankers, parmFile)
	default:
		panic(fmt.Errorf("invalid algorithm option: '%v'", optimizerAlg))
	}
}

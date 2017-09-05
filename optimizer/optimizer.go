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
func Create(optimizerAlg, parmFile string, scoreRankers []score.Ranker) Optimizer {
	switch optimizerAlg {
	case "sample":
		return NewSampleSearch(parmFile, scoreRankers)
	default:
		panic(fmt.Errorf("invalid algorithm option: '%v'", optimizerAlg))
	}
}

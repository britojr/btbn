package optimizer

import (
	"fmt"

	"github.com/britojr/btbn/score"
)

// Optimizer defines a structure optimizer algorithm
type Optimizer interface {
	Search(numSolutions, timeAvailable int) interface{}
	SetDefaultParameters()
	SetFileParameters(parms map[string]string)
}

// Create creates a structure optimizer algorithm
func Create(optimizerAlg string, scoreRankers []score.Ranker, parmFile string) Optimizer {
	switch optimizerAlg {
	case "sample":
		return newSampleSearch(scoreRankers, parmFile)
	default:
		panic(fmt.Errorf("invalid algorithm option: '%v'", optimizerAlg))
	}
}

func setParameters(alg Optimizer, parmFile string) {
	alg.SetDefaultParameters()
	parms := readParametersFile(parmFile)
	alg.SetFileParameters(parms)
}

func readParametersFile(parmFile string) map[string]string {
	return nil
}

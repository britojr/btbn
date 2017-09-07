package optimizer

import (
	"fmt"
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"

	"github.com/britojr/btbn/score"
	"github.com/britojr/utl/errchk"
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
	m := make(map[string]string)
	data, err := ioutil.ReadFile(parmFile)
	errchk.Check(err, "")
	errchk.Check(yaml.Unmarshal([]byte(data), &m), "")
	return m
}

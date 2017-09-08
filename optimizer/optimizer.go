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
	Search() *BNStructure
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

// Search applies the optimizer strategy to find the best solution
func Search(algorithm Optimizer, numSolutions, timeAvailable int) interface{} {
	var best *BNStructure
	if numSolutions <= 0 && timeAvailable <= 0 {
		numSolutions = 1
	}
	i := 0
	// remaining = available
	for {
		// timeout in remaining
		// start = tick
		current := algorithm.Search()
		// elapsed = tick - start
		// remaining -= elapsed

		if best == nil || current.Better(best) {
			best = current
		}
		// if remaining <= 0 {
		// 	break
		// }
		i++
		if numSolutions > 0 && i >= numSolutions {
			break
		}
	}
	fmt.Println(" === BEST === ")
	fmt.Printf("Score = %.6f\n", best.Score())
	return best
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

// BNStructure defines a structure solution
type BNStructure struct {
	scoreVal float64
}

// NewBNStructure creates a new structure
func NewBNStructure() *BNStructure {
	return new(BNStructure)
}

// Better returns true if this structure has a better score
func (b *BNStructure) Better(other *BNStructure) bool {
	return b.scoreVal > other.scoreVal
}

// Score returns the structure score
func (b *BNStructure) Score() float64 {
	return b.scoreVal
}

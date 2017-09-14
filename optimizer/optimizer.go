package optimizer

import (
	"io/ioutil"
	"log"
	"time"

	yaml "gopkg.in/yaml.v2"

	"github.com/britojr/btbn/score"
	"github.com/britojr/utl/errchk"
)

// Define search algorithms names
const (
	AlgSampleSearch   = "sample"    // n14
	AlgSelectedSample = "selected"  // n15
	AlgIterativeBuild = "iterative" // n16
)

// Optimizer defines a structure optimizer algorithm
type Optimizer interface {
	Search() *BNStructure
	SetDefaultParameters()
	SetFileParameters(parms map[string]string)
	PrintParameters()
}

// Create creates a structure optimizer algorithm
func Create(optimizerAlg string, scoreRankers []score.Ranker, parmFile string) (opt Optimizer) {
	switch optimizerAlg {
	case AlgSampleSearch:
		opt = NewSampleSearch(scoreRankers, parmFile)
	case AlgSelectedSample:
		opt = NewSelectSampleSearch(scoreRankers, parmFile)
	default:
		log.Panicf("invalid algorithm option: '%v'", optimizerAlg)
	}
	return
}

// Search applies the optimizer strategy to find the best solution
func Search(algorithm Optimizer, numSolutions, timeAvailable int) *BNStructure {
	var best, current *BNStructure
	if numSolutions <= 0 && timeAvailable <= 0 {
		numSolutions = 1
	}
	if timeAvailable > 0 {
		// TODO: the documentation recommends using NewTimer(d).C instead of time.After
		i := 0
		remaining := time.Duration(timeAvailable) * time.Second
		for {
			ch := make(chan *BNStructure, 1)
			start := time.Now()
			go func() {
				ch <- algorithm.Search()
			}()
			select {
			case current = <-ch:
				remaining -= time.Since(start)
			case <-time.After(remaining):
				remaining = 0
			}

			if best == nil || current.Better(best) {
				best = current
			}
			if remaining <= 0 {
				log.Printf("Time out in %v iterations\n", i)
				break
			}
			i++
			// if remaining <= 0 || (numSolutions > 0 && i >= numSolutions) {
			if numSolutions > 0 && i >= numSolutions {
				break
			}
		}
	} else {
		for i := 0; i < numSolutions; i++ {
			current := algorithm.Search()
			if best == nil || current.Better(best) {
				best = current
			}
		}
	}
	return best
}

func setParameters(alg Optimizer, parmFile string) {
	alg.SetDefaultParameters()
	if len(parmFile) > 0 {
		parms := readParametersFile(parmFile)
		alg.SetFileParameters(parms)
	}
}

func readParametersFile(parmFile string) map[string]string {
	m := make(map[string]string)
	data, err := ioutil.ReadFile(parmFile)
	errchk.Check(err, "")
	errchk.Check(yaml.Unmarshal([]byte(data), &m), "")
	return m
}

package optimizer

import (
	"fmt"
	"io/ioutil"
	"time"

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
			case <-time.After(remaining):
			}
			remaining -= time.Since(start)

			if best == nil || current.Better(best) {
				best = current
			}
			if remaining <= 0 {
				break
			}
			i++
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
	fmt.Println(" === BEST === ")
	if best != nil {
		fmt.Printf("Score = %.6f\n", best.Score())
	} else {
		fmt.Printf("Couldn't find any solution in the given time!\n")
	}
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

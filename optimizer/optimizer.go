package optimizer

import (
	"io/ioutil"
	"log"
	"time"

	yaml "gopkg.in/yaml.v2"

	"github.com/britojr/btbn/scr"
	"github.com/britojr/utl/errchk"
)

// Define nem constants
const (
	// search algorithms names
	AlgSampleSearch    = "sample"    // n14
	AlgSelectedSample  = "selected"  // n15
	AlgGuidedSearch    = "guided"    // n16
	AlgIterativeSearch = "iterative" // s16

	// file parameters
	cTreewidth       = "treewidth"
	cNumTrees        = "num_trees"
	cMutualInfo      = "mutual_info"
	cGreedy          = "greedy"
	cAstar           = "astar"
	cSearchVariation = "search_variation"
)

// Optimizer defines a structure optimizer algorithm
type Optimizer interface {
	Search() *BNStructure
	SetDefaultParameters()
	SetFileParameters(parms map[string]string)
	PrintParameters()
}

// Create creates a structure optimizer algorithm
func Create(optimizerAlg string, scoreRankers []scr.Ranker, parmFile string) (opt Optimizer) {
	switch optimizerAlg {
	case AlgSampleSearch:
		opt = NewSampleSearch(scoreRankers, parmFile)
	case AlgSelectedSample:
		opt = NewSelectSampleSearch(scoreRankers, parmFile)
	case AlgGuidedSearch:
		panic("not implemented")
	case AlgIterativeSearch:
		opt = NewIterativeSearch(scoreRankers, parmFile)
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

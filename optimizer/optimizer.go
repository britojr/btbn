package optimizer

import (
	"fmt"
	"log"
	"time"

	"github.com/britojr/btbn/scr"
)

// Define nem constants
const (
	// search algorithms names
	AlgSampleSearch    = "sample"    // n14
	AlgSelectedSample  = "selected"  // n15
	AlgGuidedSearch    = "guided"    // n16
	AlgIterativeSearch = "iterative" // s16

	// file parameters fields
	ParmTreewidth       = "treewidth"
	ParmMaxParents      = "max_parents"
	ParmNumTrees        = "num_trees"
	ParmMutualInfo      = "mutual_info"
	ParmSearchVariation = "search_variation"

	// file parameters fields options
	OpGreedy = "greedy"
	OpAstar  = "astar"
)

// Optimizer defines a structure optimizer algorithm
type Optimizer interface {
	Search() *BNStructure
	SetDefaultParameters()
	SetFileParameters(parms map[string]string)
	ValidateParameters()
	PrintParameters()
	Treewidth() int
}

var optimizerCreators = map[string]func(scr.Ranker) Optimizer{
	AlgSampleSearch:    NewSampleSearch,
	AlgSelectedSample:  NewSelectSampleSearch,
	AlgIterativeSearch: NewIterativeSearch,
}

// Create creates a structure optimizer algorithm
func Create(optimizerAlg string, scoreRanker scr.Ranker, parms map[string]string) (opt Optimizer) {
	if create, ok := optimizerCreators[optimizerAlg]; ok {
		opt = create(scoreRanker)
		setParameters(opt, parms)
		return opt
	}
	panic(fmt.Errorf("invalid algorithm option: '%v'", optimizerAlg))
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

func setParameters(alg Optimizer, parms map[string]string) {
	alg.SetDefaultParameters()
	alg.SetFileParameters(parms)
	alg.ValidateParameters()
}

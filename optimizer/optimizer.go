package optimizer

import (
	"io/ioutil"
	"log"
	"time"

	yaml "gopkg.in/yaml.v2"

	"github.com/britojr/btbn/scr"
	"github.com/britojr/utl/conv"
	"github.com/britojr/utl/errchk"
)

// Define nem constants
const (
	// search algorithms names
	AlgSampleSearch    = "sample"    // n14
	AlgSelectedSample  = "selected"  // n15
	AlgGuidedSearch    = "guided"    // n16
	AlgIterativeSearch = "iterative" // s16

	// file parameters fields
	cTreewidth       = "treewidth"
	cNumTrees        = "num_trees"
	cMutualInfo      = "mutual_info"
	cSearchVariation = "search_variation"

	// file parameters fields options
	cGreedy = "greedy"
	cAstar  = "astar"
)

// Optimizer defines a structure optimizer algorithm
type Optimizer interface {
	Search() *BNStructure
	SetDefaultParameters()
	SetFileParameters(parms map[string]string)
	PrintParameters()
	ValidateParameters()
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
	setParameters(opt, parmFile)
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
	alg.ValidateParameters()
}

func readParametersFile(parmFile string) map[string]string {
	m := make(map[string]string)
	data, err := ioutil.ReadFile(parmFile)
	errchk.Check(err, "")
	errchk.Check(yaml.Unmarshal([]byte(data), &m), "")
	return m
}

// common defines optimizer's commom/default behaviours
type common struct {
	tw           int          // treewidth
	nv           int          // number of variables
	scoreRankers []scr.Ranker // score rankers for each variable
}

func newCommon(scoreRankers []scr.Ranker, parmFile string) *common {
	s := new(common)
	s.scoreRankers = scoreRankers
	s.nv = len(s.scoreRankers)
	return s
}
func (s *common) SetDefaultParameters() {
	s.tw = 3
}
func (s *common) SetFileParameters(parms map[string]string) {
	if tw, ok := parms[cTreewidth]; ok {
		s.tw = conv.Atoi(tw)
	}
}
func (s *common) ValidateParameters() {
	if s.tw <= 0 || s.nv < s.tw+2 {
		log.Printf("n=%v, tw=%v\n", s.nv, s.tw)
		log.Panic("Invalid treewidth! Choose values such that: n >= tw+2 and tw > 0")
	}
}
func (s *common) PrintParameters() {
	log.Printf(" ========== ALGORITHM PARAMETERS ========== \n")
	log.Printf("number of variables: %v\n", s.nv)
	log.Printf("%v: %v\n", cTreewidth, s.tw)
}

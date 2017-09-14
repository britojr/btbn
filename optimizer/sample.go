package optimizer

import (
	"log"

	"github.com/britojr/btbn/ktree"
	"github.com/britojr/btbn/score"
	"github.com/britojr/utl/conv"
)

// SampleSearch implements the sampling strategy
type SampleSearch struct {
	tw           int            //treewidth
	nv           int            // number of variables
	scoreRankers []score.Ranker // score rankers for each variable
}

// NewSampleSearch creates a instance of the sample stragegy
func NewSampleSearch(scoreRankers []score.Ranker, parmFile string) *SampleSearch {
	s := new(SampleSearch)
	s.scoreRankers = scoreRankers
	s.nv = len(s.scoreRankers)
	setParameters(s, parmFile)
	s.validate()
	return s
}

// Search searchs for a network structure
func (s *SampleSearch) Search() *BNStructure {
	tk := ktree.UniformSample(s.nv, s.tw)
	bn := DAGapproximatedLearning(tk, s.scoreRankers)
	return bn
}

// SetDefaultParameters sets the defaults
func (s *SampleSearch) SetDefaultParameters() {
	// set internal variables to defined constants
	s.tw = 3
}

// SetFileParameters sets parameters from input file
func (s *SampleSearch) SetFileParameters(parms map[string]string) {
	if tw, ok := parms["treewidth"]; ok {
		s.tw = conv.Atoi(tw)
	}
}

// Validate validates parameters
func (s *SampleSearch) validate() {
	if s.tw <= 0 || s.nv < s.tw+2 {
		log.Printf("n=%v, tw=%v\n", s.nv, s.tw)
		log.Panic("Invalid treewidth! Choose values such that: n >= tw+2 and tw > 0")
	}
}

// PrintParameters prints the algorithm's current parameters
func (s *SampleSearch) PrintParameters() {
	log.Printf(" ========== ALGORITHM PARAMETERS ========== \n")
	log.Printf("number of variables: %v\n", s.nv)
	log.Printf("treewidth: %v\n", s.tw)
	log.Printf(" ------------------------------------------ \n")
}

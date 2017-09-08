package optimizer

import (
	"log"
	"time"

	"github.com/britojr/btbn/score"
	"github.com/britojr/utl/conv"
)

// SampleSearch implements the sampling strategy
type sampleSearch struct {
	tw           int            //treewidth
	nv           int            // number of variables
	scoreRankers []score.Ranker // score rankers for each variable
}

// NewSampleSearch creates a instance of the sample stragegy
func newSampleSearch(scoreRankers []score.Ranker, parmFile string) Optimizer {
	s := new(sampleSearch)
	s.scoreRankers = scoreRankers
	s.nv = len(s.scoreRankers)
	setParameters(s, parmFile)
	s.validate()
	return s
}

// Search search for a network structure
func (s *sampleSearch) Search() *BNStructure {
	// ctree.Sample(s.nv, s.tw) sample a ktree
	// dag.LearnFromKTree(ctree, s.scoreRankers) // returns a bnstruct (parentset list and score)
	bn := NewBNStructure()
	bn.scoreVal = -1
	time.Sleep(1 * time.Second)
	log.Println("a search")
	return bn
}

// SetDefaultParameters sets the defaults
func (s *sampleSearch) SetDefaultParameters() {
	// set internal variables to defined constants
	s.tw = 3
}

// SetFileParameters sets parameters from input file
func (s *sampleSearch) SetFileParameters(parms map[string]string) {
	if tw, ok := parms["treewidth"]; ok {
		s.tw = conv.Atoi(tw)
	}
}

// Validate validates parameters
func (s *sampleSearch) validate() {
	if s.tw <= 0 || s.nv < s.tw+2 {
		log.Printf("n=%v, tw=%v\n", s.nv, s.tw)
		log.Panic("Invalid treewidth! Choose values such that: n >= tw+2 and tw > 0")
	}
}

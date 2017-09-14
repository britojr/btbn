package optimizer

import (
	"log"

	"github.com/britojr/btbn/ktree"
	"github.com/britojr/btbn/score"
	"github.com/britojr/tcc/codec"
	"github.com/britojr/utl/conv"
)

// SelectSampleSearch implements the select sampling strategy
type SelectSampleSearch struct {
	scoreRankers []score.Ranker // score rankers for each variable
	nv           int            // number of variables
	tw           int            // treewidth
	prevCodes    []*codec.Code  // previously accepted codes
	bestTk       *ktree.Ktree   // currently best scoring tree
	numTrees     int            // number of ktrees to sample before start learning DAG
	tkList       []*ktree.Ktree // list of accepted ktrees
}

// NewSelectSampleSearch creates a instance of the sample stragegy
func NewSelectSampleSearch(scoreRankers []score.Ranker, parmFile string) *SelectSampleSearch {
	s := new(SelectSampleSearch)
	s.scoreRankers = scoreRankers
	s.nv = len(s.scoreRankers)
	setParameters(s, parmFile)
	s.validate()
	return s
}

// Search searchs for a network structure
func (s *SelectSampleSearch) Search() *BNStructure {
	if len(s.tkList) == 0 {
		s.selectKTrees()
	}
	tk := s.tkList[0]
	s.tkList = s.tkList[1:]
	bn := DAGapproximatedLearning(tk, s.scoreRankers)
	return bn
}

// SetDefaultParameters sets the defaults
func (s *SelectSampleSearch) SetDefaultParameters() {
	// set internal variables to defined constants
	s.tw = 3
	s.numTrees = 1
}

// SetFileParameters sets parameters from input file
func (s *SelectSampleSearch) SetFileParameters(parms map[string]string) {
	if tw, ok := parms["treewidth"]; ok {
		s.tw = conv.Atoi(tw)
	}
	if numTrees, ok := parms["num_trees"]; ok {
		s.numTrees = conv.Atoi(numTrees)
	}
}

// Validate validates parameters
func (s *SelectSampleSearch) validate() {
	if s.tw <= 0 || s.nv < s.tw+2 {
		log.Printf("n=%v, tw=%v\n", s.nv, s.tw)
		log.Panic("Invalid treewidth! Choose values such that: n >= tw+2 and tw > 0")
	}
}

// PrintParameters prints the algorithm's current parameters
func (s *SelectSampleSearch) PrintParameters() {
	log.Printf(" ========== ALGORITHM PARAMETERS ========== \n")
	log.Printf("number of variables: %v\n", s.nv)
	log.Printf("treewidth: %v\n", s.tw)
	log.Printf(" ------------------------------------------ \n")
}

// selectKTrees samples and selects a given number of ktrees
func (s *SelectSampleSearch) selectKTrees() {
	s.tkList = make([]*ktree.Ktree, s.numTrees)
	for i := range s.tkList {
		s.tkList[i] = ktree.UniformSample(s.nv, s.tw)
	}
}

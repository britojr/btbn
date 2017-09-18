package optimizer

import (
	"log"

	"github.com/britojr/btbn/scr"
	"github.com/britojr/utl/conv"
)

// IterativeSearch implements the dag iterative building strategy
type IterativeSearch struct {
	scoreRankers []scr.Ranker // score rankers for each variable
	nv           int          // number of variables
	tw           int          // treewidth

	algVariant string // algorithm variant

	prevCliques map[string]int // previously sampled initial cliques
}

// NewIterativeSearch creates a instance of the iterative stragegy
func NewIterativeSearch(scoreRankers []scr.Ranker, parmFile string) *IterativeSearch {
	s := new(IterativeSearch)
	s.scoreRankers = scoreRankers
	s.nv = len(s.scoreRankers)
	setParameters(s, parmFile)
	s.validate()
	s.prevCliques = make(map[string]int)
	return s
}

// Search searchs for a network structure
func (s *IterativeSearch) Search() *BNStructure {
	// ord = sample order
	// bn = learn exact DAG for initial k+1 clique (ord[:tw+1])
	// bn = search from starting node (bn, ord[tw+1:])
	return nil
}

// SetDefaultParameters sets the defaults
func (s *IterativeSearch) SetDefaultParameters() {
	// set internal variables to defined constants
	s.tw = 3
	s.algVariant = "greedy"
}

// SetFileParameters sets parameters from input file
func (s *IterativeSearch) SetFileParameters(parms map[string]string) {
	if tw, ok := parms["treewidth"]; ok {
		s.tw = conv.Atoi(tw)
	}
	if algVariant, ok := parms["alg_variant"]; ok {
		s.algVariant = algVariant
	}
}

// Validate validates parameters
func (s *IterativeSearch) validate() {
	if s.tw <= 0 || s.nv < s.tw+2 {
		log.Printf("n=%v, tw=%v\n", s.nv, s.tw)
		log.Panic("Invalid treewidth! Choose values such that: n >= tw+2 and tw > 0")
	}
	if !(s.algVariant == "greedy" || s.algVariant == "astar") {
		log.Panicf("Invalid algorithm variant option: '%v'", s.algVariant)
	}
}

// PrintParameters prints the algorithm's current parameters
func (s *IterativeSearch) PrintParameters() {
	log.Printf(" ========== ALGORITHM PARAMETERS ========== \n")
	log.Printf("number of variables: %v\n", s.nv)
	log.Printf("treewidth: %v\n", s.tw)
	log.Printf("search type: '%v'\n", s.algVariant)
	log.Printf(" ------------------------------------------ \n")
}

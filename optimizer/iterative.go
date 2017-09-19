package optimizer

import (
	"fmt"
	"log"
	"math/rand"
	"sort"

	"github.com/britojr/btbn/ktree"
	"github.com/britojr/btbn/scr"
	"github.com/britojr/btbn/varset"
	"github.com/britojr/utl/conv"
)

// IterativeSearch implements the dag iterative building strategy
type IterativeSearch struct {
	scoreRankers []scr.Ranker // score rankers for each variable
	nv           int          // number of variables
	tw           int          // treewidth

	searchVariation string // search variation

	prevCliques map[string]struct{} // previously sampled initial cliques
}

// NewIterativeSearch creates a instance of the iterative stragegy
func NewIterativeSearch(scoreRankers []scr.Ranker, parmFile string) *IterativeSearch {
	s := new(IterativeSearch)
	s.scoreRankers = scoreRankers
	s.nv = len(s.scoreRankers)
	setParameters(s, parmFile)
	s.validate()
	s.prevCliques = make(map[string]struct{})
	return s
}

// Search searchs for a network structure
func (s *IterativeSearch) Search() *BNStructure {
	ord := s.sampleOrder()
	bn := s.getInitialDAG(ord[:s.tw+1])
	switch s.searchVariation {
	case cGreedy:
		s.greedySearch(bn, ord)
	case cAstar:
		panic("not implemented")
	default:
		log.Panicf("invalid search variation: '%v'", s.searchVariation)
	}
	return bn
}

// SetDefaultParameters sets the defaults
func (s *IterativeSearch) SetDefaultParameters() {
	// set internal variables to defined constants
	s.tw = 3
	s.searchVariation = cGreedy
}

// SetFileParameters sets parameters from input file
func (s *IterativeSearch) SetFileParameters(parms map[string]string) {
	if tw, ok := parms[cTreewidth]; ok {
		s.tw = conv.Atoi(tw)
	}
	if searchVariation, ok := parms[cSearchVariation]; ok {
		s.searchVariation = searchVariation
	}
}

// Validate validates parameters
func (s *IterativeSearch) validate() {
	if s.tw <= 0 || s.nv < s.tw+2 {
		log.Printf("n=%v, tw=%v\n", s.nv, s.tw)
		log.Panic("Invalid treewidth! Choose values such that: n >= tw+2 and tw > 0")
	}
	if !(s.searchVariation == cGreedy || s.searchVariation == cAstar) {
		log.Panicf("Invalid algorithm variant option: '%v'", s.searchVariation)
	}
}

// PrintParameters prints the algorithm's current parameters
func (s *IterativeSearch) PrintParameters() {
	log.Printf(" ========== ALGORITHM PARAMETERS ========== \n")
	log.Printf("number of variables: %v\n", s.nv)
	log.Printf("%v: %v\n", cTreewidth, s.tw)
	log.Printf("%v: '%v'\n", cSearchVariation, s.searchVariation)
	log.Printf(" ------------------------------------------ \n")
}

func (s *IterativeSearch) sampleOrder() []int {
	r := rand.New(rand.NewSource(seed()))
	for {
		ord := r.Perm(s.nv)
		// key := varset.New(s.nv).SetInts(ord[:s.tw+1]).DumpHashString()
		sort.Ints(ord[:s.tw+1])
		key := fmt.Sprint(ord[:s.tw+1])
		if _, ok := s.prevCliques[key]; !ok {
			s.prevCliques[key] = struct{}{}
			return ord
		}
	}
}

func (s *IterativeSearch) getInitialDAG(vars []int) *BNStructure {
	// TODO: replace this for an exact method
	bn := DAGapproximatedLearning(ktree.New(vars, -1, -1), s.scoreRankers)
	return bn
}

func (s *IterativeSearch) greedySearch(bn *BNStructure, ord []int) *BNStructure {
	// clqs := []varset.Varset{varset.New(s.nv)}
	clqs := make([]varset.Varset, 0, s.nv-s.tw)
	clqs = append(clqs, varset.New(s.nv))
	for _, v := range ord[:s.tw+1] {
		clqs[0].Set(v)
	}
	ord = ord[s.tw+1:]
	for len(ord) > 0 {
		v := ord[0]
		bestPs, bestScr := s.scoreRankers[v].BestInLim(clqs[0], s.tw)
		for _, clq := range clqs[1:] {
			ps, scr := s.scoreRankers[v].BestInLim(clq, s.tw)
			if scr > bestScr {
				bestScr, bestPs = scr, ps
			}
		}
		bn.SetParents(v, bestPs, bestScr)
		clqs = append(clqs, bestPs.Clone().Set(v))
		ord = ord[1:]
	}
	return bn
}

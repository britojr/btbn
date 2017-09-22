package optimizer

import (
	"fmt"
	"log"
	"math/rand"
	"sort"

	"github.com/britojr/btbn/ktree"
	"github.com/britojr/btbn/scr"
	"github.com/britojr/btbn/varset"
)

// IterativeSearch implements the dag iterative building strategy
type IterativeSearch struct {
	*common                             // common variables and methods
	searchVariation string              // search variation
	prevCliques     map[string]struct{} // previously sampled initial cliques
}

// NewIterativeSearch creates an instance of the iterative stragegy
func NewIterativeSearch(scoreRankers []scr.Ranker) Optimizer {
	s := &IterativeSearch{common: newCommon(scoreRankers)}
	s.prevCliques = make(map[string]struct{})
	return s
}

// Search searchs for a network structure
func (s *IterativeSearch) Search() *BNStructure {
	ord := s.sampleOrder()
	bn := s.getInitialDAG(ord[:s.tw+1])
	switch s.searchVariation {
	case OpGreedy:
		s.greedySearch(bn, ord)
	case OpAstar:
		panic("not implemented")
	default:
		log.Panicf("invalid search variation: '%v'", s.searchVariation)
	}
	return bn
}

// SetDefaultParameters sets parameters to default values
func (s *IterativeSearch) SetDefaultParameters() {
	s.common.SetDefaultParameters()
	s.searchVariation = OpGreedy
}

// SetFileParameters sets parameters from input file
func (s *IterativeSearch) SetFileParameters(parms map[string]string) {
	s.common.SetFileParameters(parms)
	if searchVariation, ok := parms[ParmSearchVariation]; ok {
		s.searchVariation = searchVariation
	}
}

// ValidateParameters validates internal parameters
func (s *IterativeSearch) ValidateParameters() {
	s.common.ValidateParameters()
	if !(s.searchVariation == OpGreedy || s.searchVariation == OpAstar) {
		log.Panicf("Invalid algorithm variant option: '%v'", s.searchVariation)
	}
}

// PrintParameters prints the algorithm's current parameters
func (s *IterativeSearch) PrintParameters() {
	s.common.PrintParameters()
	log.Printf("%v: '%v'\n", ParmSearchVariation, s.searchVariation)
}

// sampleOrder samples a permutation of variables
// rejecting repeated k+1 initial variables that already occured in previous samples
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
	bestBn := DAGapproximatedLearning(ktree.New(vars, -1, -1), s.scoreRankers)
	for i := 0; i < 50; i++ {
		currBn := DAGapproximatedLearning(ktree.New(vars, -1, -1), s.scoreRankers)
		if currBn.Better(bestBn) {
			bestBn = currBn
		}
	}
	return bestBn
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

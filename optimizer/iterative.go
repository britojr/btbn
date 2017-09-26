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

	// to use on astar
	order []int
	hval  []float64
}

// NewIterativeSearch creates an instance of the iterative stragegy
func NewIterativeSearch(scoreRanker scr.Ranker) Optimizer {
	s := &IterativeSearch{common: newCommon(scoreRanker)}
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
		s.astarSearch(bn, ord)
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
	bestBn := DAGapproximatedLearning(ktree.New(vars, -1, -1), s.scoreRanker)
	for i := 0; i < 50; i++ {
		currBn := DAGapproximatedLearning(ktree.New(vars, -1, -1), s.scoreRanker)
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
		// the list has cliques of size k+1, hence v should have at most k parents
		// in order to form another k+1 clique
		bestPs, bestScr := s.scoreRanker.BestInLim(v, clqs[0], s.tw)
		for _, clq := range clqs[1:] {
			ps, scr := s.scoreRanker.BestInLim(v, clq, s.tw)
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

type problemState struct {
	next  int           // index of next variable to be included
	clqs  [][]int       // list of cliques
	v     int           // variable that was assigned to reach this state
	pset  varset.Varset // parent set assigned to variable to reach this state
	pscor float64       // partial score for the variables currently included
}

type searchNode struct {
	state  *problemState
	parent *searchNode
	score  float64 // accumulated solution score
}

func (s *IterativeSearch) astarSearch(bn *BNStructure, ord []int) *BNStructure {
	state := s.getStartState(bn, ord)
	pq := &scr.RecordSlice{scr.NewRecord(state.pscor+s.heuristic(state), &searchNode{state, nil, state.pscor})}
	for pq.Len() > 0 {
		nd := pq.Pop().(*scr.Record).Data().(*searchNode)
		if s.isGoalState(nd.state) {
			return s.makeSolution(bn, nd)
		}
		for _, succ := range s.stateSuccessors(nd.state) {
			ch := &searchNode{succ, nd, nd.score + succ.pscor}
			pq.Push(scr.NewRecord(ch.score+s.heuristic(succ), ch))
		}
	}
	return nil
}

func (s *IterativeSearch) makeSolution(bn *BNStructure, nd *searchNode) *BNStructure {
	for nd.parent != nil {
		bn.SetParents(nd.state.v, nd.state.pset, nd.state.pscor)
		nd = nd.parent
	}
	return bn
}

func (s *IterativeSearch) getStartState(bn *BNStructure, ord []int) *problemState {
	s.order = ord
	s.hval = make([]float64, len(ord))
	restric := varset.New(s.nv)
	for _, u := range ord {
		restric.Set(u)
	}
	for i := len(s.hval) - 1; i > s.tw; i-- {
		restric.Clear(ord[i])
		_, scor := s.scoreRanker.BestIn(ord[i], restric)
		s.hval[i] = scor
		if i < len(s.hval)-1 {
			s.hval[i] += s.hval[i+1]
		}
	}
	st := new(problemState)
	st.v = -1
	st.pset = nil
	st.next = s.tw + 1
	st.clqs = append(st.clqs, ord[:s.tw+1])
	for _, v := range ord[:s.tw+1] {
		st.pscor += bn.LocalScore(v)
	}
	return st
}

func (s *IterativeSearch) isGoalState(ps *problemState) bool {
	return ps.next >= len(s.order)
}

func (s *IterativeSearch) heuristic(ps *problemState) float64 {
	if s.isGoalState(ps) {
		return 0
	}
	return s.hval[ps.next]
}

func (s *IterativeSearch) stateSuccessors(ps *problemState) (succ []*problemState) {
	v := s.order[ps.next]
	clique := ps.clqs[0]
	succ = append(succ, s.successorClique(clique, v, ps)...)
	for _, clq := range ps.clqs[1:] {
		succ = append(succ, s.successorClique(clq[1:], v, ps)...)
	}
	return
}

func (s *IterativeSearch) successorClique(clique []int, v int, ps *problemState) (succ []*problemState) {
	allset := varset.New(s.nv).Set(v)
	for _, u := range clique {
		allset.Set(u)
	}
	for i, u := range clique {
		clq := append([]int{v}, clique[:i]...)
		clq = append(clq, clique[i+1:]...)
		restric := allset.Clone().Clear(u)
		pset, pscor := s.scoreRanker.BestIn(v, restric)
		succ = append(succ, &problemState{
			next:  ps.next + 1,
			clqs:  append(ps.clqs, clq),
			v:     v,
			pset:  pset,
			pscor: pscor,
		})
	}
	return
}

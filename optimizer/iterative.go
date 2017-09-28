package optimizer

import (
	"container/heap"
	"fmt"
	"log"
	"math/rand"

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
// rejecting repeated orders (k+1 to forward) that already occured on previous samples
func (s *IterativeSearch) sampleOrder() []int {
	r := rand.New(rand.NewSource(seed()))
	for {
		ord := r.Perm(s.nv)
		key := fmt.Sprint(ord[s.tw+1:])
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
	s.order = ord
	s.computeHeuristic()
	state := s.getStartState(bn)
	rs := []*scr.Record{scr.NewRecord(-(state.pscor + s.heuristic(state)), &searchNode{state, nil, state.pscor})}
	pq := scr.NewRecordHeap(&rs, func(i, j int) bool { return rs[i].Score() < rs[j].Score() })
	heap.Init(pq)
	cpush, cpop := 0, 0
	fmt.Printf("in: %v\n", rs[0].Score())
	for pq.Len() > 0 {
		// uncomment:
		// nd := heap.Pop(pq).(*scr.Record).Data().(*searchNode)
		// remove:
		aux := heap.Pop(pq).(*scr.Record)
		// fmt.Printf("pop: %v | %v\n", aux.Score(), aux.Data().(*searchNode).state.clqs)
		cpop++
		nd := aux.Data().(*searchNode)
		// remove up
		if s.isGoalState(nd.state) {
			fmt.Printf(">>> tot push:%v, pop:%v\n", cpush, cpop)
			return s.makeSolution(bn, nd)
		}
		for _, succ := range s.stateSuccessors(nd.state) {
			ch := &searchNode{succ, nd, nd.score + succ.pscor}
			// uncomment:
			// heap.Push(pq, scr.NewRecord(-(ch.score+s.heuristic(succ)), ch))
			// remove:
			aux := scr.NewRecord(-(ch.score + s.heuristic(succ)), ch)
			// fmt.Printf("  push: %v | %v\n", aux.Score(), aux.Data().(*searchNode).state.clqs)
			heap.Push(pq, aux)
			cpush++
			// remove up
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

func (s *IterativeSearch) computeHeuristic() {
	s.hval = make([]float64, len(s.order))
	restric := varset.New(s.nv)
	for _, u := range s.order {
		restric.Set(u)
	}
	for i := len(s.hval) - 1; i > s.tw; i-- {
		restric.Clear(s.order[i])
		_, scor := s.scoreRanker.BestIn(s.order[i], restric)
		s.hval[i] = scor
		if i < len(s.hval)-1 {
			s.hval[i] += s.hval[i+1]
		}
	}
}

func (s *IterativeSearch) getStartState(bn *BNStructure) *problemState {
	st := new(problemState)
	st.v = -1
	st.pset = nil
	st.next = s.tw + 1
	st.clqs = append(st.clqs, s.order[:s.tw+1])
	for _, v := range s.order[:s.tw+1] {
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
	succ = append(succ, s.successorClique(clique, []int{v}, ps)...)
	for _, clq := range ps.clqs[1:] {
		succ = append(succ, s.successorClique(clq[1:], []int{v, clq[0]}, ps)...)
	}
	return
}

func (s *IterativeSearch) successorClique(clique []int, pref []int, ps *problemState) (succ []*problemState) {
	v := pref[0]
	allset := varset.New(s.nv)
	for _, u := range clique {
		allset.Set(u)
	}
	for _, u := range pref[1:] {
		allset.Set(u)
	}
	for i, u := range clique {
		clq := append(pref, clique[:i]...)
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

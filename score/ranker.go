package score

import (
	"fmt"
	"sort"

	"github.com/britojr/btbn/varset"
)

// Ranker defines a list of best scores for a given variable
type Ranker interface {
	BestIn(restrictSet varset.Varset) (parents varset.Varset, localScore float64)
	ScoreOf(parents varset.Varset) float64
}

// CreateRankers creates array of rankers, one for each variable
func CreateRankers(cache *Cache, maxPa int) []Ranker {
	rs := []Ranker(nil)
	for i := 0; i < cache.Nvar(); i++ {
		rs = append(rs, NewListRanker(i, cache, maxPa))
	}
	return rs
}

// ListRanker trivial implementation of score ranker
type ListRanker struct {
	varIndex  int
	scoreList []varsetScore
	scoreMap  map[string]float64
}

// NewListRanker creates new listRanker
func NewListRanker(varIndex int, cache *Cache, maxPa int) *ListRanker {
	m := &ListRanker{}
	m.varIndex = varIndex
	m.scoreMap = cache.Scores(varIndex)
	m.scoreList = make([]varsetScore, 0, len(m.scoreMap))
	for s, scor := range m.scoreMap {
		pset := varset.New(cache.Nvar())
		pset.LoadHashString(s)
		if maxPa <= 0 || pset.Count() <= maxPa {
			m.scoreList = append(m.scoreList, varsetScore{scor, pset})
		}
	}
	sort.Sort(varsetScores(m.scoreList))
	return m
}

// BestIn finds the highest scoring parent set that is contained in the given restriction set
func (m *ListRanker) BestIn(restric varset.Varset) (parents varset.Varset, scr float64) {
	if len(m.scoreList) == 0 {
		panic(fmt.Errorf("Score list is empty"))
	}
	for _, v := range m.scoreList {
		if restric.IsSuperSet(v.vars) {
			return v.vars, v.scor
		}
	}
	panic(fmt.Errorf("Can't find score for variable %v with restriction %v", m.varIndex, restric))
}

// ScoreOf returns the score of a given set of parents
func (m *ListRanker) ScoreOf(parents varset.Varset) float64 {
	return m.scoreMap[parents.DumpHashString()]
}

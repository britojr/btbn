package scr

import (
	"fmt"

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
	scoreList []Record
	scoreMap  map[string]float64
}

// NewListRanker creates new listRanker
func NewListRanker(varIndex int, cache *Cache, maxPa int) *ListRanker {
	m := &ListRanker{}
	m.varIndex = varIndex
	m.scoreMap = cache.Scores(varIndex)
	m.scoreList = make([]Record, 0, len(m.scoreMap))
	for s, scor := range m.scoreMap {
		pset := varset.New(cache.Nvar())
		pset.LoadHashString(s)
		if maxPa <= 0 || pset.Count() <= maxPa {
			m.scoreList = append(m.scoreList, Record{scor, pset})
		}
	}
	SortRecord(m.scoreList)
	return m
}

// BestIn finds the highest scoring parent set that is contained in the given restriction set
func (m *ListRanker) BestIn(restric varset.Varset) (parents varset.Varset, scr float64) {
	if len(m.scoreList) == 0 {
		panic(fmt.Errorf("Score list is empty"))
	}
	for _, v := range m.scoreList {
		parents = v.data.(varset.Varset)
		if restric.IsSuperSet(parents) {
			return parents, v.score
		}
	}
	panic(fmt.Errorf("Can't find score for variable %v with restriction %v", m.varIndex, restric))
}

// ScoreOf returns the score of a given set of parents
func (m *ListRanker) ScoreOf(parents varset.Varset) float64 {
	return m.scoreMap[parents.DumpHashString()]
}

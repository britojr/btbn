package score

import (
	"fmt"
	"sort"

	"github.com/britojr/utl/conv"
	"github.com/willf/bitset"
)

// Ranker defines a list of best scores for a given variable
type Ranker interface {
	BestIn(restrictSet *bitset.BitSet) (parents *bitset.BitSet, localScore float64)
}

// CreateRankers creates array of rankers, one for each variable
func CreateRankers(cache *Cache, maxPa int) []Ranker {
	rs := []Ranker(nil)
	for i := 0; i < cache.Nvar(); i++ {
		rs = append(rs, NewListRanker(i, cache, maxPa))
	}
	return rs
}

// listRanker trivial implementation of score ranker
type listRanker struct {
	varIndex  int
	maxPa     int
	scoreList []varsetScore
}

// NewListRanker creates new mapRanker
func NewListRanker(varIndex int, cache *Cache, maxPa int) Ranker {
	m := &listRanker{}
	m.varIndex = varIndex
	m.maxPa = maxPa
	scoreMap := cache.Scores(varIndex)
	m.scoreList = make([]varsetScore, 0, len(scoreMap))
	for k, v := range scoreMap {
		m.scoreList = append(m.scoreList, varsetScore{v, conv.AtoBs(k)})
	}
	sort.Sort(varsetScores(m.scoreList))
	return m
}

// BestIn finds the highest scoring parent set that is contained in the given restriction set
func (m *listRanker) BestIn(restric *bitset.BitSet) (parents *bitset.BitSet, scr float64) {
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

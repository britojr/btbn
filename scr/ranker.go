package scr

import (
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/britojr/btbn/varset"
	"github.com/britojr/utl/conv"
	"github.com/britojr/utl/ioutl"
)

// Ranker defines methods to find best score for a given variable
type Ranker interface {
	BestIn(v int, restric varset.Varset) (parents varset.Varset, localScore float64)
	BestInLim(v int, restric varset.Varset, maxPa int) (parents varset.Varset, localScore float64)
	ScoreOf(v int, parents varset.Varset) float64
	Size() int
	SaveSubSet(fname string, vs []int)
}

// CreateRanker creates a ranker for variables given in cache
func CreateRanker(cache *Cache, maxPa int) Ranker {
	r := new(rankList)
	for i := 0; i < cache.Nvar(); i++ {
		r.vars = append(r.vars, newVarRanker(i, cache, maxPa))
	}
	return r
}

type rankList struct {
	vars []*varRanker
}

// Size returns the number of variables in the ranker
func (r *rankList) Size() int {
	return len(r.vars)
}

// ScoreOf returns the score of a family
func (r *rankList) ScoreOf(v int, parents varset.Varset) float64 {
	return r.vars[v].scoreOf(parents)
}

// BestIn finds the highest scoring parent set that is contained in the given restriction set
func (r *rankList) BestIn(v int, restric varset.Varset) (parents varset.Varset, scr float64) {
	return r.vars[v].bestIn(restric, restric.Count())
}

// BestInLim finds the highest scoring parent set that is contained in the given restriction set
// and respects a maximum size of the parent set
func (r *rankList) BestInLim(v int, restric varset.Varset, maxPa int) (parents varset.Varset, scr float64) {
	return r.vars[v].bestIn(restric, maxPa)
}

func (r *rankList) SaveSubSet(fname string, vs []int) {
	// log.Printf("clique (%v): %v\n", fname, vs)
	w := ioutl.CreateFile(fname)
	vset := varset.New(r.Size())
	vset.SetInts(vs)
	for _, v := range vs {
		if r.vars[v].varIndex != v {
			log.Panicf("ranker: wrong index! %v != %v\n", r.vars[v].varIndex, v)
		}
		fmt.Fprintf(w, "VAR %v\n", v)
		for _, rec := range r.vars[v].scoreList {
			pset := rec.Data().(varset.Varset)
			if vset.IsSuperSet(pset) {
				fmt.Fprintf(w, "%v %v\n", rec.Score(), strings.Join(conv.Sitoa(pset.DumpAsInts()), " "))
			}
		}
		fmt.Fprintln(w)
	}
	w.Close()
}

type varRanker struct {
	varIndex  int
	scoreList []*Record
	scoreMap  map[string]float64
}

func newVarRanker(v int, cache *Cache, maxPa int) *varRanker {
	r := new(varRanker)
	r.varIndex = v
	r.scoreMap = cache.Scores(v)
	r.scoreList = make([]*Record, 0, len(r.scoreMap))
	for s, score := range r.scoreMap {
		pset := varset.New(cache.Nvar())
		pset.LoadHashString(s)
		if maxPa <= 0 || pset.Count() <= maxPa {
			r.scoreList = append(r.scoreList, NewRecord(score, pset))
		}
	}
	sort.Slice(r.scoreList, func(i, j int) bool {
		return r.scoreList[i].Score() > r.scoreList[j].Score()
	})
	return r
}

func (r *varRanker) bestIn(restric varset.Varset, maxPa int) (parents varset.Varset, score float64) {
	if len(r.scoreList) == 0 {
		panic(fmt.Errorf("Score list is empty"))
	}
	for _, v := range r.scoreList {
		parents = v.Data().(varset.Varset)
		if restric.IsSuperSet(parents) && parents.Count() <= maxPa {
			// if restric.IsSuperSet(parents) {
			return parents, v.score
		}
	}
	panic(fmt.Errorf("Can't find score for variable %v with restriction %v", r.varIndex, restric))
}

func (r *varRanker) scoreOf(parents varset.Varset) float64 {
	return r.scoreMap[parents.DumpHashString()]
}

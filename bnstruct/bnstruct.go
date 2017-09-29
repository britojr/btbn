package bnstruct

import (
	"fmt"
	"math"

	"github.com/britojr/btbn/varset"
)

// BNStruct defines a bn structure as a list of parent sets
type BNStruct struct {
	totScore    float64
	initialized int
	nodes       []*node
}

type node struct {
	locScore float64
	parents  varset.Varset
}

// New creates a new bn structure
func New(nvars int) *BNStruct {
	b := new(BNStruct)
	b.nodes = make([]*node, nvars)
	return b
}

// Better returns true if this structure has a better score
func (b *BNStruct) Better(other *BNStruct) bool {
	// return (other == nil) || b.Score() > other.Score()
	if other == nil {
		return true
	}
	if b.initialized == other.initialized {
		return b.totScore > other.totScore
	}
	return b.initialized > other.initialized
}

// Score returns the structure total score
// if the structure has at least one variable uninitialized, returns math.Inf(-1)
func (b *BNStruct) Score() float64 {
	if b.Size() != b.initialized {
		return math.Inf(-1)
	}
	return b.totScore
}

// Size returns the number of variables in the network
func (b *BNStruct) Size() int {
	return len(b.nodes)
}

// SetParents assigns the parent set for a variable and its respective local score
func (b *BNStruct) SetParents(v int, parents varset.Varset, localScore float64) {
	if b.nodes[v] == nil {
		b.nodes[v] = new(node)
		b.initialized++
	}
	b.totScore -= b.nodes[v].locScore
	b.nodes[v].locScore = localScore
	b.nodes[v].parents = parents
	b.totScore += localScore
}

// LocalScore returns the score of variable v's family
// if the score is undefined for v, returns math.Inf(-1)
func (b *BNStruct) LocalScore(v int) float64 {
	if b.nodes[v] == nil {
		return math.Inf(-1)
	}
	return b.nodes[v].locScore
}

// Parents returns the parent set of a variable
func (b *BNStruct) Parents(v int) varset.Varset {
	if b.nodes[v] == nil {
		return nil
	}
	return b.nodes[v].parents
}

func (b *BNStruct) String() string {
	s := fmt.Sprintf("{size: %v\n", b.Size())
	for i := 0; i < b.Size(); i++ {
		s += fmt.Sprintf("\t%v: {s(%v), p(%v)}\n", i, b.LocalScore(i), b.Parents(i))
	}
	s += fmt.Sprintf("total: %v}\n", b.Score())
	return s
}

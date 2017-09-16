package optimizer

import (
	"math"

	"github.com/britojr/btbn/varset"
)

// BNStructure defines a structure solution
type BNStructure struct {
	totScore    float64
	initialized int
	nodes       []*node
}

type node struct {
	locScore float64
	parents  varset.Varset
}

// NewBNStructure creates a new structure
func NewBNStructure(nvars int) *BNStructure {
	b := new(BNStructure)
	b.nodes = make([]*node, nvars)
	return b
}

// Better returns true if this structure has a better score
func (b *BNStructure) Better(other *BNStructure) bool {
	return (other == nil) || b.Score() > other.Score()
}

// Score returns the structure score
func (b *BNStructure) Score() float64 {
	if b.Size() != b.initialized {
		return -math.MaxFloat64
	}
	return b.totScore
}

// Size returns the number of variables in the network
func (b *BNStructure) Size() int {
	return len(b.nodes)
}

// SetParents assigns the parent set for a variable and its respective local score
func (b *BNStructure) SetParents(v int, parents varset.Varset, localScore float64) {
	if b.nodes[v] == nil {
		b.nodes[v] = new(node)
		b.initialized++
	}
	b.totScore -= b.nodes[v].locScore
	b.nodes[v].locScore = localScore
	b.nodes[v].parents = parents
	b.totScore += localScore
}

// LocalScore returns the local score of a family
func (b *BNStructure) LocalScore(v int) float64 {
	if b.nodes[v] == nil {
		return -math.MaxFloat64
	}
	return b.nodes[v].locScore
}

// Parents returns the parent set of a variable
func (b *BNStructure) Parents(v int) varset.Varset {
	if b.nodes[v] == nil {
		return nil
	}
	return b.nodes[v].parents
}

package optimizer

import "github.com/britojr/btbn/varset"

// BNStructure defines a structure solution
type BNStructure struct {
	scoreVal float64
}

// NewBNStructure creates a new structure
func NewBNStructure() *BNStructure {
	return new(BNStructure)
}

// Better returns true if this structure has a better score
func (b *BNStructure) Better(other *BNStructure) bool {
	return (other == nil) || b.scoreVal > other.scoreVal
}

// Score returns the structure score
func (b *BNStructure) Score() float64 {
	return b.scoreVal
}

// SetParents assigns the parent set for a variable and its respective local score
func (b *BNStructure) SetParents(v int, parents varset.Varset, localScore float64) {
	panic("not implemented")
}

// LocalScore returns the local score of a family
func (b *BNStructure) LocalScore(v int) float64 {
	panic("not implemented")
}

package optimizer

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

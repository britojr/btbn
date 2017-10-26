package model

import "github.com/britojr/btbn/factor"

// CTree defines a structure in clique tree format
// a CTree is a way to group the potentials of the model according to its cliques
// the potentials assossiated with each clique are pointers to the same factors present in the model
type CTree struct{}

func (c *CTree) NCliques() int {
	panic("ctree: not implemented")
}
func (c *CTree) RootID() int {
	panic("ctree: not implemented")
}
func (c *CTree) Neighbors(v int) []int {
	panic("ctree: not implemented")
}
func (c *CTree) VarIn(v int) []int {
	panic("ctree: not implemented")
}
func (c *CTree) VarOut(v int) []int {
	panic("ctree: not implemented")
}

func (c *CTree) Potentials() []*factor.Factor {
	panic("ctree: not implemented")
}
func (c *CTree) SetPotentials([]*factor.Factor) {
	panic("ctree: not implemented")
}
func ToCTree(m Model) *CTree {
	panic("ctree: not implemented")
}

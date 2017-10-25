package model

import (
	"github.com/britojr/btbn/factor"
	"github.com/britojr/btbn/vars"
)

// CTree defines a structure in clique tree format
// a CTree is a way to group the potentials of the model according to its cliques
// the potentials assossiated with each clique are pointers to the same factors present in the model
type CTree struct{}

func (c *CTree) NCliques() int {
	panic("ctree: not implemented")
}
func (c *CTree) Potentials() []*factor.Factor {
	panic("ctree: not implemented")
}
func (c *CTree) Neighb(v int) []int {
	panic("ctree: not implemented")
}
func (c *CTree) VarIn(v int) *vars.Var {
	panic("ctree: not implemented")
}
func (c *CTree) VarOut(v int) *vars.Var {
	panic("ctree: not implemented")
}

func ToCTree(m Model) *CTree {
	panic("ctree: not implemented")
}

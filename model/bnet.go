package model

import (
	"github.com/britojr/btbn/factor"
	"github.com/britojr/btbn/vars"
)

// Model defines a probabilistic graphical model
// type Model interface {
// 	Node(*vars.Var) *BNode
// }

// BNet defines a Bayesian network model
type BNet struct {
	nodes map[*vars.Var]*BNode
}

// BNode defines a BN node
type BNode struct {
	vx  *vars.Var
	cpt *factor.Factor
}

// NewBNet creates a new BNet model
func NewBNet() *BNet {
	b := new(BNet)
	b.nodes = make(map[*vars.Var]*BNode)
	return b
}

// ToCTree return a ctree for this bnet
func (b *BNet) ToCTree() *CTree {
	panic("not implemented")
}

// Node return the respective node of a var
func (b *BNet) Node(v *vars.Var) *BNode {
	return b.nodes[v]
}

// Variable return node pivot variable
func (nd *BNode) Variable() *vars.Var {
	return nd.vx
}

// Potential return node potential
func (nd *BNode) Potential() *factor.Factor {
	return nd.cpt
}

// SetPotential set node potential
func (nd *BNode) SetPotential(p *factor.Factor) {
	nd.cpt = p
}

// func (b BNet) SetFamily(varx *vars.Var, cpt *factor.Factor) {
// 	b.nodes[varx] = &bnode{cpt}
// }

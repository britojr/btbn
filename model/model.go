package model

import (
	"github.com/britojr/btbn/factor"
	"github.com/britojr/btbn/vars"
)

// Model defines a probabilistic graphical model
type Model interface {
	Copy() Model
	// Plus(Model) Model
	// SetParameters(Model) Model
	// Normalize() Model
	// CTree() *CTree
	// SetCTree(*CTree)
	// Variables() vars.VarList
	// SetCPT(*factor.Factor)
}

type bnet struct {
	nodes map[*vars.Var]*bnode
}
type bnode struct {
	// varx *vars.Var
	cpt *factor.Factor
}

// NewBNet creates a new bnet model
func NewBNet() Model {
	b := new(bnet)
	b.nodes = make(map[*vars.Var]*bnode)
	return b
}

func (b bnet) SetFamily(varx *vars.Var, cpt *factor.Factor) {
	b.nodes[varx] = &bnode{cpt}
}

func (b bnet) CTree() *CTree {
	panic("model: not implemented")
}

func (b bnet) SetCTree(*CTree) {
	panic("model: not implemented")
}

func (b bnet) Copy() Model {
	panic("model: not implemented")
}
func (b bnet) Plus(Model) Model {
	panic("model: not implemented")
}
func (b bnet) SetParameters(Model) Model {
	panic("model: not implemented")
}
func (b bnet) Normalize() Model {
	panic("model: not implemented")
}

// func (b bnet) Variables() vars.VarList {
// 	panic("model: not implemented")
// }
//
// func (b bnet) SetCPT(cpt *factor.Factor) {
// 	panic("model: not implemented")
// }

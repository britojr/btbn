package model

import (
	"github.com/britojr/btbn/vars"
	"github.com/britojr/kbn/factor"
)

// Model defines a probabilistic graphical model
type Model interface {
	// Variables() vars.VarList
	// SetCPT(*factor.Factor)
	Copy() Model
	Plus(Model) Model
	SetParameters(Model) Model
	Normalize() Model
}

type bnet struct {
	nodes map[*vars.Var]*bnode
}
type bnode struct {
	// varx *vars.Var
	cpt *factor.Factor
}

// NewBNet creates a new bnet model
// func NewBNet() Model {
func NewBNet() *bnet {
	b := new(bnet)
	b.nodes = make(map[*vars.Var]*bnode)
	return b
}

func (b bnet) SetFamily(varx *vars.Var, cpt *factor.Factor) {
	b.nodes[varx] = &bnode{cpt}
}

func (b bnet) SetCTree() {
	panic("model: not implemented")
}

// func (b bnet) Variables() vars.VarList {
// 	panic("model: not implemented")
// }
//
// func (b bnet) SetCPT(cpt *factor.Factor) {
// 	panic("model: not implemented")
// }

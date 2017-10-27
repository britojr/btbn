package inference

import (
	"github.com/britojr/btbn/dataset"
	"github.com/britojr/btbn/factor"
	"github.com/britojr/btbn/model"
	"github.com/gonum/floats"
)

// InfAlg defines an inference algorithm
type InfAlg interface {
	Run(dataset.Evidence) float64
	CTNodes() []*model.CTNode
	CalibPotential(nd *model.CTNode) *factor.Factor
	UpdatedModel() *model.BNet
}

type cTCalib struct {
	ct                *model.CTree
	bn                *model.BNet
	size              int
	initPot, calibPot map[*model.CTNode]*factor.Factor

	// auxiliar for message passing, send to parent and receive from parent
	send, receive map[*model.CTNode]*factor.Factor
	// axiliar to reduce (memoize) number of factor multiplications
	prev, post map[*model.CTNode][]*factor.Factor
}

// NewCTreeCalibration creates a new clique tree calibration runner
func NewCTreeCalibration(bn *model.BNet) InfAlg {
	c := new(cTCalib)
	c.bn = bn
	ct := bn.ToCTree()
	c.ct = ct
	c.size = ct.Len()

	// initialize slices to be used on calibration
	c.initPot = make(map[*model.CTNode]*factor.Factor)
	c.calibPot = make(map[*model.CTNode]*factor.Factor)
	c.send = make(map[*model.CTNode]*factor.Factor)
	c.receive = make(map[*model.CTNode]*factor.Factor)
	c.prev = make(map[*model.CTNode][]*factor.Factor)
	c.post = make(map[*model.CTNode][]*factor.Factor)
	return c
}

// UpdatedModel updates model's parameters based on the internal ctree
func (c *cTCalib) UpdatedModel() *model.BNet {
	// TODO: check if it needs to calibrate without evidence first
	// c.Run(map[int]int{})
	for v, nd := range c.ct.Families() {
		bnd := c.bn.Node(v)
		p := nd.Potential().Copy().Marginalize(bnd.Potential().Variables()...)
		p.Normalize(bnd.Variable())
		bnd.SetPotential(p)
	}
	return c.bn
}

// Nodes returns reference for ctree nodes
func (c *cTCalib) CTNodes() []*model.CTNode {
	return c.ct.Nodes()
}

// CalibPotential returns calibrated potential of a node
func (c *cTCalib) CalibPotential(nd *model.CTNode) *factor.Factor {
	return c.calibPot[nd]
}

func (c *cTCalib) Run(e dataset.Evidence) float64 {
	c.applyEvidence(e)
	c.upDownCalibration()
	// after applying evidence and calibratin
	// the sum of any potential is probability of evidence
	return floats.Sum(c.calibPot[c.ct.Root()].Values())
}

// applyEvidence initialize the potentials with a copy of the original potentials
// applyed the given evidence
func (c *cTCalib) applyEvidence(e dataset.Evidence) {
	for _, nd := range c.ct.Nodes() {
		c.initPot[nd] = nd.Potential().Copy().Reduce(e)
	}
}

// upDownCalibration runs two-passage message passing clique tree calibration
// by the end, every node should have the joint distribution of its respective clique variables
func (c *cTCalib) upDownCalibration() {
	// -------------------------------------------------------------------------
	// send[nd] contains the message the node nd sends up to its parent
	// receive[nd] contains the message the node nd receives from his parent
	// -------------------------------------------------------------------------
	// post[nd][j] contains the product of every message that node nd
	// received from its (j+1)th child up to the last child
	// prev[nd][j] contains the product of node nd initial potential and
	// every message that node nd received from its fist child to the (j-1)th child
	// So the message to be sent from node nd to its jth chilnd
	// will be the product of prev[nd][j] and post[nd][j]
	// -------------------------------------------------------------------------

	root := c.ct.Root()
	c.upwardmessage(root, nil)
	c.downwardmessage(nil, root)
}

func (c *cTCalib) upwardmessage(v, pa *model.CTNode) {
	children := v.Children()
	c.prev[v] = make([]*factor.Factor, len(children)+1)
	c.prev[v][0] = c.initPot[v]
	for i, ch := range children {
		c.upwardmessage(ch, v)
		c.prev[v][i+1] = c.send[ch].Copy().Times(c.prev[v][i])
	}
	if pa != nil {
		c.send[v] = c.prev[v][len(c.prev[v])-1].Copy().Marginalize(pa.Variables()...)
	}
}

func (c *cTCalib) downwardmessage(pa, v *model.CTNode) {
	children := v.Children()
	c.calibPot[v] = c.prev[v][len(c.prev[v])-1].Copy()
	if pa != nil {
		c.calibPot[v].Times(c.receive[v])
	}
	// if v is a leaf, nothing more to do
	if len(children) == 0 {
		return
	}

	c.post[v] = make([]*factor.Factor, len(children))
	i := len(c.post[v]) - 1
	c.post[v][i] = c.receive[v]
	i--
	for ; i >= 0; i-- {
		ch := children[i+1]
		c.post[v][i] = c.send[ch].Copy()
		if c.post[v][i+1] != nil {
			c.post[v][i].Times(c.post[v][i+1])
		}
	}

	for i, ch := range children {
		msg := c.prev[v][i].Copy()
		if c.post[v][i] != nil {
			msg.Times(c.post[v][i])
		}
		c.receive[ch] = msg.Marginalize(ch.Variables()...)
		c.downwardmessage(v, ch)
	}
}

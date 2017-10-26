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
	CalibPotList() []*factor.Factor
	SetOrigPotList([]*factor.Factor)
	OrigPotList() []*factor.Factor
	SetModelParms(m model.Model) model.Model
}

type cTCalib struct {
	ct                                *model.CTree
	size                              int
	initPot, calibPot, calibPotSepSet []*factor.Factor

	// auxiliar for message passing, send to parent and receive from parent
	send, receive []*factor.Factor
	// axiliar to reduce (memoize) number of factor multiplications
	prev, post [][]*factor.Factor
}

// NewCTreeCalibration creates a new clique tree calibration runner
func NewCTreeCalibration(ct *model.CTree) InfAlg {
	c := new(cTCalib)
	c.ct = ct
	c.size = ct.NCliques()

	// initialize slices to be used on calibration
	c.initPot = make([]*factor.Factor, c.size)
	c.calibPot = make([]*factor.Factor, c.size)
	c.calibPotSepSet = make([]*factor.Factor, c.size)
	c.send = make([]*factor.Factor, c.size)
	c.receive = make([]*factor.Factor, c.size)
	c.prev = make([][]*factor.Factor, c.size)
	c.post = make([][]*factor.Factor, c.size)
	return c
}

// SetModelParms updates model's parameters based on the internal ctree
func (c *cTCalib) SetModelParms(m model.Model) model.Model {
	panic("inference: not implemented")
}

// OrigPotList returns reference for ctree original parameters
func (c *cTCalib) OrigPotList() []*factor.Factor {
	return c.ct.Potentials()
}

// SetOrigPotList updates internat ctree parameters
func (c *cTCalib) SetOrigPotList(ps []*factor.Factor) {
	c.ct.SetPotentials(ps)
}

// CalibPotList returns calibrated potentials
func (c *cTCalib) CalibPotList() []*factor.Factor {
	return c.calibPot
}

func (c *cTCalib) Run(e dataset.Evidence) float64 {
	c.applyEvidence(e)
	c.upDownCalibration()
	// after applying evidence and calibratin
	// the sum of any potential is probability of evidence
	return floats.Sum(c.calibPot[0].Values())
}

// applyEvidence initialize the potentials with a copy of the original potentials
// applyed the given evidence
func (c *cTCalib) applyEvidence(e dataset.Evidence) {
	for i, p := range c.ct.Potentials() {
		c.initPot[i] = p.Copy().Reduce(e)
	}
}

// upDownCalibration runs two-passage message passing clique tree calibration
// by the end, every node should have the joint distribution of its respective clique variables
func (c *cTCalib) upDownCalibration() {
	// -------------------------------------------------------------------------
	// send[i] contains the message the ith node sends up to its parent
	// receive[i] contains the message the ith node receives from his parent
	// -------------------------------------------------------------------------
	// post[i][j] contains the product of every message that node i received
	// from its j+1 children to the last children
	// prev[i][j] contains the product of node i initial potential and
	// every message that node i received from its fist children to the j-1 children
	// So the message to be sent from i to j will be the product of prev and post
	// -------------------------------------------------------------------------

	root := c.ct.RootID()
	c.upwardmessage(root, -1)
	c.downwardmessage(-1, root)
}

func (c *cTCalib) upwardmessage(v, pa int) {
	neighbors := c.ct.Neighbors(v)
	c.prev[v] = make([]*factor.Factor, 1, len(neighbors)+1)
	c.prev[v][0] = c.initPot[v]
	for _, ne := range neighbors {
		if ne != pa {
			c.upwardmessage(ne, v)
			c.prev[v] = append(c.prev[v], c.send[ne].TimesNew(c.prev[v][len(c.prev[v])-1]))
		}
	}
	if pa != -1 {
		c.send[v] = c.prev[v][len(c.prev[v])-1].SumOutIDNew(c.ct.VarIn(v)...)
	}
}

func (c *cTCalib) downwardmessage(pa, v int) {
	neighbors := c.ct.Neighbors(v)
	c.calibPot[v] = c.prev[v][len(c.prev[v])-1]
	n := len(neighbors)
	if pa != -1 {
		c.calibPot[v].Times(c.receive[v])
		n--
		// calculate calibrated sepset
		c.calibPotSepSet[v] = c.calibPot[v].SumOutIDNew(c.ct.VarIn(v)...)
	}
	if len(neighbors) == 1 && pa != -1 {
		return
	}

	c.post[v] = make([]*factor.Factor, n)
	i := len(c.post[v]) - 1
	c.post[v][i] = c.receive[v]
	i--
	for k := len(neighbors) - 1; k >= 0 && i >= 0; k-- {
		ch := neighbors[k]
		if ch == pa {
			continue
		}
		c.post[v][i] = c.send[ch]
		if c.post[v][i+1] != nil {
			c.post[v][i].Times(c.post[v][i+1])
		}
		i--
	}

	k := 0
	for _, ch := range neighbors {
		if ch == pa {
			continue
		}
		msg := c.prev[v][k].Copy()
		if c.post[v][k] != nil {
			msg.Times(c.post[v][k])
		}
		c.receive[ch] = msg.SumOutID(c.ct.VarOut(ch)...)
		c.downwardmessage(v, ch)
		k++
	}
}

package inference

import (
	"github.com/britojr/btbn/dataset"
	"github.com/britojr/btbn/factor"
	"github.com/britojr/btbn/model"
)

// InfAlg defines an inference algorithm
type InfAlg interface {
	Run(dataset.Evidence) float64
	CalibPotList() []*factor.Factor
	SetOrigPotList([]*factor.Factor)
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
func NewCTreeCalibration(ct *model.CTree) *cTCalib {
	c := new(cTCalib)
	c.ct = ct
	c.size = ct.NCliques()
	c.calibPot = make([]*factor.Factor, c.size)
	return c
}

// SetModelParms sets model's parameters based on the internal ctree
func (c *cTCalib) SetModelParms(m model.Model) model.Model {
	panic("inference: not implemented")
}

// SetOrigPotList updates internat ctree parameters
func (c *cTCalib) SetOrigPotList(ps []*factor.Factor) {
	panic("inference: not implemented")
}

// CalibPotList returns calibrated potentials
func (c *cTCalib) CalibPotList() []*factor.Factor {
	panic("inference: not implemented")
}

func (c *cTCalib) Run(e dataset.Evidence) float64 {
	c.applyEvidence(e)
	c.upDownCalibration()
	panic("inference: not implemented")
}

// applyEvidence initialize the potentials with a copy of the original potentials
// applyed the given evidence
func (c *cTCalib) applyEvidence(e dataset.Evidence) {
	for i, p := range c.ct.Potentials() {
		c.initPot[i] = p.Copy().Reduce(e)
	}
}

func (c *cTCalib) upDownCalibration() {
	// -------------------------------------------------------------------------
	// send[i] contains the message the ith node sends up to its parent
	// receive[i] contains the message the ith node receives from his parent
	// -------------------------------------------------------------------------
	c.send = make([]*factor.Factor, c.size)
	c.receive = make([]*factor.Factor, c.size)
	// -------------------------------------------------------------------------
	// post[i][j] contains the product of every message that node i received
	// from its j+1 children to the last children
	// prev[i][j] contains the product of node i initial potential and
	// every message that node i received from its fist children to the j-1 children
	// So the message to be sent from i to j will be the product of prev and post
	// -------------------------------------------------------------------------
	c.prev = make([][]*factor.Factor, c.size)
	c.post = make([][]*factor.Factor, c.size)

	c.calibPot = make([]*factor.Factor, c.size)
	c.calibPotSepSet = make([]*factor.Factor, c.size)
	root := 0

	c.upwardmessage(root, -1)
	c.downwardmessage(-1, root)
}

func (c *cTCalib) upwardmessage(v, pa int) {
	neighb := c.ct.Neighb(v)
	c.prev[v] = make([]*factor.Factor, 1, len(neighb)+1)
	c.prev[v][0] = c.initPot[v]
	for _, ne := range neighb {
		if ne != pa {
			c.upwardmessage(ne, v)
			c.prev[v] = append(c.prev[v], c.send[ne].TimesNew(c.prev[v][len(c.prev[v])-1]))
		}
	}
	if pa != -1 {
		c.send[v] = c.prev[v][len(c.prev[v])-1].SumOutNew(c.ct.VarIn(v))
	}
}

func (c *cTCalib) downwardmessage(pa, v int) {
	neighb := c.ct.Neighb(v)
	c.calibPot[v] = c.prev[v][len(c.prev[v])-1]
	n := len(neighb)
	if pa != -1 {
		c.calibPot[v].Times(c.receive[v])
		n--
		// calculate calibrated sepset
		c.calibPotSepSet[v] = c.calibPot[v].SumOutNew(c.ct.VarIn(v))
	}
	if len(neighb) == 1 && pa != -1 {
		return
	}

	c.post[v] = make([]*factor.Factor, n)
	i := len(c.post[v]) - 1
	c.post[v][i] = c.receive[v]
	i--
	for k := len(neighb) - 1; k >= 0 && i >= 0; k-- {
		ch := neighb[k]
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
	for _, ch := range neighb {
		if ch == pa {
			continue
		}
		msg := c.prev[v][k].Copy()
		if c.post[v][k] != nil {
			msg.Times(c.post[v][k])
		}
		c.receive[ch] = msg.SumOut(c.ct.VarOut(ch))
		c.downwardmessage(v, ch)
		k++
	}
}

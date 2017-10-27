// Package emlearner implements Expectation-Maximization algorithm for parameter estimation
package emlearner

import (
	"log"
	"math"

	"github.com/britojr/btbn/dataset"
	"github.com/britojr/btbn/factor"
	"github.com/britojr/btbn/inference"
	"github.com/britojr/btbn/model"
	"github.com/britojr/utl/conv"
)

// property map options
const (
	ParmMaxIters  = "em_max_iters" // maximum number of iterations
	ParmThreshold = "em_threshold" // minimum improvement threshold
)

// default properties
const (
	cMaxIters  = 5
	cThreshold = 1e-1
)

// EMLearner implements Expectation-Maximization algorithm
type EMLearner interface {
	SetProperties(props map[string]string)
	Run(model.Model, dataset.EvidenceSet) (model.Model, float64)
}

// implementation of EMLearner
type emAlg struct {
	maxIters  int     // max number of em iterations
	threshold float64 // minimum improvement threshold
	nIters    int     // number of iterations of current alg
}

// New creates a new EMLearner
func New() EMLearner {
	e := new(emAlg)
	// set defaults
	e.maxIters = cMaxIters
	e.threshold = cThreshold
	return e
}

func (e *emAlg) SetProperties(props map[string]string) {
	// set properties
	if maxIters, ok := props[ParmMaxIters]; ok {
		e.maxIters = conv.Atoi(maxIters)
	}
	if threshold, ok := props[ParmThreshold]; ok {
		e.threshold = conv.Atof(threshold)
	}
	// validate properties
	if e.maxIters <= 0 {
		log.Panicf("emlearner: max iterations (%v) must be > 0", e.maxIters)
	}
	if e.threshold <= 0 {
		log.Panicf("emlearner: convergence threshold (%v) must be > 0", e.threshold)
	}
}

// start defines a starting point for model's parameters
func (e *emAlg) start(infalg inference.InfAlg, evset dataset.EvidenceSet) {
	// TODO: add a non-trivial em (re)start policy
	// for now, just randomly starts
	for _, nd := range infalg.CTNodes() {
		nd.Potential().RandomDistribute()
	}
}

// Run runs EM until convergence or max iteration number is reached
func (e *emAlg) Run(m model.Model, evset dataset.EvidenceSet) (model.Model, float64) {
	log.Printf("emlearner: start\n")
	infalg := inference.NewCTreeCalibration(model.ToCTree(m))
	e.nIters = 0
	e.start(infalg, evset)
	var llant, llnew float64
	for {
		llnew = e.runStep(infalg, evset)
		e.nIters++
		if llant != 0 && (e.nIters >= e.maxIters || (math.Abs((llnew-llant)/llant) < e.threshold)) {
			break
		}
		log.Printf("\temlearner: diff=%v\n", math.Abs((llnew-llant)/llant))
		llant = llnew
	}
	log.Printf("emlearner: iterations=%v\n", e.nIters)
	return infalg.SetModelParms(m), llnew
}

// runStep runs expectation and maximization steps
// returning the loglikelihood of the model with new parameters
func (e *emAlg) runStep(infalg inference.InfAlg, evset dataset.EvidenceSet) float64 {
	// copy of parameters to hold the sufficient statistics
	count := make(map[*model.CTNode]*factor.Factor)
	var ll float64
	// expecttation step
	for _, evid := range evset.Observations() {
		// evid is a map of var to state
		evLkhood := infalg.Run(evid)
		ll += math.Log(evLkhood)

		// acumulate sufficient statistics in the copy of parameters
		for _, nd := range infalg.CTNodes() {
			if p, ok := count[nd]; ok {
				p.Plus(infalg.CalibPotential(nd))
			} else {
				count[nd] = infalg.CalibPotential(nd).Copy()
			}
		}
	}

	// maximization step
	// updates parameters
	for nd, p := range count {
		p.Normalize()
		nd.SetPotential(p)
	}

	// updates loglikelihood of optimized model
	// m.SetLoglikelihood(ds, ll)
	return ll
}

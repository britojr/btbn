// Package emlearner implements Expectation-Maximization algorithm for parameter estimation
package emlearner

import (
	"log"
	"math"

	"github.com/britojr/btbn/dataset"
	"github.com/britojr/btbn/inference"
	"github.com/britojr/btbn/model"
)

// EMLearner implements Expectation-Maximization algorithm
type EMLearner interface {
	SetProperties(props map[string]string)
	Run(model.Model, dataset.EvidenceSet) (loglikelihood float64)
}

// implementation of EMLearner
type emAlg struct {
	maxIters  int     // max number of em iterations
	threshold float64 // minimum improvement threshold
	nIters    int     // number of iterations of current alg
}

func (e *emAlg) SetProperties(props map[string]string) {
	panic("emlearner: not implemented")
}

func (e *emAlg) Run(m model.Model, evset dataset.EvidenceSet) (model.Model, float64) {
	log.Printf("emlearner: start\n")
	e.nIters = 0
	infalg := e.start(m, evset)
	var llant, llnew float64
	for {
		m, llnew = e.runStep(infalg, evset)
		e.nIters++
		if llant != 0 && (e.nIters >= e.maxIters || (math.Abs((llnew-llant)/llant) < e.threshold)) {
			break
		}
		log.Printf("\temlearner: diff=%v\n", math.Abs((llnew-llant)/llant))
		llant = llnew
	}
	log.Printf("emlearner: iterations=%v\n", e.nIters)
	return m, llnew
}

func (e *emAlg) start(m model.Model, evset dataset.EvidenceSet) inference.InfAlg {
	// define a starting point for model's parameters
	// create an inference alg with the model
	panic("emlearner: not implemented")
}

func (e *emAlg) runStep(infalg inference.InfAlg, evset dataset.EvidenceSet) (model.Model, float64) {
	// expecttation step
	// use a copy of the model to hold the sufficient statistics
	var count, m model.Model = nil, nil
	var ll float64
	for _, evid := range evset.Observations() {
		// evid is a map of var to state
		infalg.SetEvidence(evid)
		evidLikelihood := infalg.Run()
		ll += math.Log(evidLikelihood)

		// acumulates sufficient statistics on the copy model
		m = infalg.Model()
		if count == nil {
			count = m.Copy()
		} else {
			count.Plus(m)
		}
	}

	// maximization step
	// updates parameters
	m.SetParameters(count.Normalize())
	// updates loglikelihood of optimized model
	// m.SetLoglikelihood(ds, ll)
	return m, ll
}

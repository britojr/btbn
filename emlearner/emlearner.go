// Package emlearner implements Expectation-Maximization algorithm for parameter estimation
package emlearner

import (
	"log"
	"math"

	"github.com/britojr/btbn/dataset"
	"github.com/britojr/btbn/factor"
	"github.com/britojr/btbn/inference"
	"github.com/britojr/btbn/model"
)

// EMLearner implements Expectation-Maximization algorithm
type EMLearner interface {
	SetProperties(props map[string]string)
	Run(model.BNet, dataset.EvidenceSet) (loglikelihood float64)
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

func (e *emAlg) Run(bn model.BNet, evset dataset.EvidenceSet) (model.BNet, float64) {
	log.Printf("emlearner: start\n")
	e.nIters = 0
	infalg := e.start(bn, evset)
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
	return infalg.BNet(), llnew
}

func (e *emAlg) start(bn model.BNet, evset dataset.EvidenceSet) inference.CTreeCalibration {
	// define a starting point for model's parameters
	// create an inference alg with the model
	panic("emlearner: not implemented")
}

func (e *emAlg) runStep(infalg inference.CTreeCalibration, evset dataset.EvidenceSet) float64 {
	// sufficient statistics for each node
	count := make(map[int]factor.Factor)

	// runs expecttation step
	var ll float64
	for _, evid := range evset.Observations() {
		// evid is a map of var to state
		infalg.SetEvidence(evid)
		evidLikelihood := infalg.Run()
		ll += math.Log(evidLikelihood)

		// updates sufficient statistics for each node
		bn := infalg.BNet()
		for _, v := range bn.Variables() {
			frac := infalg.FamilyBelief(v)
			if _, ok := count[v.ID()]; !ok {
				count[v.ID()].Plus(frac)
			} else {
				count[v.ID()] = frac
			}
		}
	}

	// runs maximization step
	// updates parameters
	bn := infalg.BNet()
	for _, v := range bn.Variables() {
		count[v.ID()].Normalize(v)
		bn.SetCPT(count[v.ID()])
	}
	// updates loglikelihood of optimized model
	// m.SetLoglikelihood(ds, ll)
	return ll
}

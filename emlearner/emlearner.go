// Package emlearner implements Expectation-Maximization algorithm for parameter estimation
package emlearner

import (
	"log"
	"math"
)

// EMLearner implements Expectation-Maximization algorithm
type EMLearner interface {
	SetProperties(props map[string]string)
	Run(model, dataset) (loglikelihood float64)
}

// temporary interfaces
type model interface{}
type evidenceset interface{}
type infAlg interface {
	Model() model
}

type emAlg struct {
	maxIters  int     // max number of em iterations
	threshold float64 // minimum improvement threshold
	nIters    int     // number of iterations of current alg
}

func (e *emAlg) SetProperties(props map[string]string) {
	panic("emlearner: not implemented")
}

func (e *emAlg) Run(m model, ds evidenceset) (model, float64) {
	log.Printf("emlearner: start\n")
	e.nIters = 0
	infalg := e.start(m, ds)
	var llant, llnew float64
	for {
		llnew = e.runStep(infalg, ds)
		e.nIters++
		if llant != 0 && (e.nIters >= e.maxIters || (math.Abs((llnew-llant)/llant) < e.threshold)) {
			break
		}
		log.Printf("\temlearner: diff=%v\n", math.Abs((llnew-llant)/llant))
		llant = llnew
	}
	log.Printf("emlearner: iterations=%v\n", e.nIters)
	return infalg.Model(), llnew
}

func (e *emAlg) start(m model, ds evidenceset) infAlg {
	// define a starting point for model's parameters
	// create an inference alg with the model
	panic("emlearner: not implemented")
}

func (e *emAlg) runStep(infalg infAlg, ds evidenceset) float64 {
	// runs expecttation step
	// runs maximization step

	// sufficient statistics for each node
	count := make(map[int]factor)

	var ll float64
	for _, evid := range ds.Maps() {
		// evid is a map of var to state
		infalg.SetEvidence(evid)
		evidLikelihood = m.Propagate()
		ll += math.Log(evidLikelihood)

		// updates sufficient statistics for each node
		for i, node := range infalg.Model().Nodes() {
			frac := infalg.ComputeFamilyBelief(node)
			if count[i] {
				count[i].Plus(frac)
			} else {
				count[i] = frac
			}
		}
	}

	// updates parameters
	for i, node := range infalg.Model().Nodes() {
		cpt := count[i]
		cpt.Normalize(node.Variable())
		Node.setCpt(cpt)
	}

	// updates loglikelihood of optimized model
	m.setLoglikelihood(ds, ll)
	panic("emlearner: not implemented")
}

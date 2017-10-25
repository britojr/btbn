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

func (e *emAlg) start(infalg inference.InfAlg, evset dataset.EvidenceSet) {
	// define a starting point for model's parameters
	// create an inference alg with the model
	panic("emlearner: not implemented")
}

func (e *emAlg) runStep(infalg inference.InfAlg, evset dataset.EvidenceSet) float64 {
	// expecttation step
	// to hold the sufficient statistics
	count := make([]*factor.Factor, len(infalg.CalibPotList()))
	var ll float64
	for _, evid := range evset.Observations() {
		// evid is a map of var to state
		evLkhood := infalg.Run(evid)
		ll += math.Log(evLkhood)

		// acumulates sufficient statistics on the copy of parameters
		ps := infalg.CalibPotList()
		for i, p := range ps {
			if count[i] == nil {
				count[i] = p.Copy()
			} else {
				count[i].Plus(p)
			}
		}
	}

	// maximization step
	// updates parameters
	for i := range count {
		count[i].Normalize()
	}
	infalg.SetOrigPotList(count)

	// updates loglikelihood of optimized model
	// m.SetLoglikelihood(ds, ll)
	return ll
}

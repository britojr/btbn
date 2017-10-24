package inference

import (
	"github.com/britojr/btbn/dataset"
	"github.com/britojr/btbn/model"
)

// InfAlg defines an inference algorithm
type InfAlg interface {
	Model() model.Model
	SetEvidence(dataset.Evidence)
	Run() float64
	// FamilyBelief(*vars.Var) *factor.Factor
}

type cTreeCalibration struct{}

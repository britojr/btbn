package inference

import (
	"github.com/britojr/btbn/dataset"
	"github.com/britojr/btbn/factor"
	"github.com/britojr/btbn/model"
	"github.com/britojr/btbn/variable"
)

// CTreeCalibration ..
type CTreeCalibration interface {
	BNet() model.BNet
	SetEvidence(dataset.Evidence)
	Run() float64
	FamilyBelief(variable.Var) factor.Factor
}

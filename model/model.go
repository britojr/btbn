package model

import (
	"github.com/britojr/btbn/factor"
	"github.com/britojr/btbn/variable"
)

// BNet ..
type BNet interface {
	Variables() []variable.Var
	SetCPT(factor.Factor)
}

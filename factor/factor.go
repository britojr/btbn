package factor

import "github.com/britojr/btbn/variable"

// Factor ..
type Factor interface {
	Plus(Factor)
	Normalize(...variable.Var)
}

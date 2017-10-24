package vars

import "fmt"

const (
	// DefaultNState default number of states for a variable
	DefaultNState = 2
)

// Var defines variable type
type Var struct {
	id, nstate int
}

// New creates a new variable
func New(id, nstate int) *Var {
	return &Var{id, nstate}
}

// ID variable's id
func (v Var) ID() int {
	return v.id
}

// NState variable's number of states
func (v Var) NState() int {
	return v.nstate
}

func (v Var) String() string {
	return fmt.Sprintf("x%v[%v]", v.id, v.nstate)
}

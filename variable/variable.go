package variable

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

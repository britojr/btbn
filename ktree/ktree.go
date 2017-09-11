package ktree

// Ktree defines a ktree structure
type Ktree struct {
}

// Variables returns the variables of this node
func (t *Ktree) Variables() []int {
	panic("not implemented")
}

// Children returns the children of this node
func (t *Ktree) Children() []*Ktree {
	panic("not implemented")
}

// VarIn returns the variable that was included in this node relative to the parent node
func (t *Ktree) VarIn() int {
	panic("not implemented")
}

// VarOut returns the variable that was removed from this node relative to the parent node
func (t *Ktree) VarOut() int {
	panic("not implemented")
}

// UniformSample uniformly samples a ktree
func UniformSample(n, k int) *Ktree {
	panic("not implemented")
}

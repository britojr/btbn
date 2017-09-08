package ctree

// Ctree defines a ktree in clique tree form
type Ctree struct {
}

// Variables returns the clique variables of this node
func (c *Ctree) Variables() []int {
	panic("not implemented")
}

// Children returns the children of this node
func (c *Ctree) Children() []*Ctree {
	panic("not implemented")
}

// VarIn returns variable that was included relative to parent
func (c *Ctree) VarIn() int {
	panic("not implemented")
}

// VarOut returns variable that was removed relative to parent
func (c *Ctree) VarOut() int {
	panic("not implemented")
}

// SepSet returns the intersection with parent
func (c *Ctree) SepSet() []int {
	panic("not implemented")
}

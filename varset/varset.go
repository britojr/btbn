package varset

import "fmt"

// Varset defines a variable set type
type Varset interface {
	DumpHashString() string
	LoadHashString(s string) Varset
	DumpAsInts() []int
	SetInts(is []int) Varset
	ClearInts(is []int) Varset
	IsSuperSet(other Varset) bool
	Set(i int) Varset
	Clear(i int) Varset
	Count() int
	Equal(other Varset) bool
	Clone() Varset
}

const (
	varsetTypeDefault = "uibset" // default varset implementation
	// varset type implementations
	typeUibset = "uibset" // implementation using unsigned ints
	typeBibset = "bibset" // implementation using big ints
)

var varsetCreators = map[string]func(int) Varset{
	typeUibset: newUibset,
	// typeBibset:  newBibset,
}

// Create creates a varset of defined type
func Create(varsetType string, size int) Varset {
	if create, ok := varsetCreators[varsetType]; ok {
		return create(size)
	}
	panic(fmt.Errorf("invalid option: '%v'", varsetType))
}

// New creates new varset
func New(size int) Varset {
	return Create(varsetTypeDefault, size)
}

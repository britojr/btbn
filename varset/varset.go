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

var varsetTypeDefault = "uibset"

// var varsetTypeDefault = "bigbset"

// Create creates a varset of defined type
func Create(varsetType string, size int) Varset {
	switch varsetType {
	case "uibset":
		return newUibset(size)
	default:
		panic(fmt.Errorf("invalid option: '%v'", varsetType))
	}
}

// New creates new varset
func New(size int) Varset {
	return Create(varsetTypeDefault, size)
}

package varset

import (
	"github.com/britojr/utl/errchk"
	"github.com/willf/bitset"
)

type uibset struct {
	*bitset.BitSet
}

func newUibset(size int) Varset {
	return &uibset{bitset.New(uint(size))}
}

func (b *uibset) DumpHashString() string {
	mbytes, err := b.MarshalBinary()
	errchk.Check(err, "")
	return string(mbytes)
}

func (b *uibset) DumpAsInts() []int {
	// return b.DumpAsInts()
	s := make([]int, 0, b.Count())
	for i, ok := b.NextSet(0); ok; i, ok = b.NextSet(i + 1) {
		s = append(s, int(i))
	}
	return s
}

func (b *uibset) LoadHashString(s string) Varset {
	b.UnmarshalBinary([]byte(s))
	return b
}

func (b *uibset) IsSuperSet(other Varset) bool {
	return b.BitSet.IsSuperSet(other.(*uibset).BitSet)
}

func (b *uibset) Set(i int) Varset {
	b.BitSet.Set(uint(i))
	return b
}
func (b *uibset) Clear(i int) Varset {
	b.BitSet.Clear(uint(i))
	return b
}

func (b *uibset) Count() int {
	return int(b.BitSet.Count())
}

func (b *uibset) Equal(other Varset) bool {
	return b.BitSet.Equal(other.(*uibset).BitSet)
}

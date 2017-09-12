package varset

import "github.com/willf/bitset"

type uibset struct {
	*bitset.BitSet
}

func newUibset(size int) Varset {
	return &uibset{bitset.New(uint(size))}
}

func (b *uibset) DumpAsString() string {
	return b.DumpAsBits()
}

func (b *uibset) DumpAsInts() []int {
	// return b.DumpAsInts()
	s := make([]int, 0, b.Count())
	for i, ok := b.NextSet(0); ok; i, ok = b.NextSet(i + 1) {
		s = append(s, int(i))
	}
	return s
}

func (b *uibset) SetFromString(s string) {
	for i, v := range s {
		if v == '1' {
			b.BitSet.Set(uint(i))
		}
	}
}

func (b *uibset) IsSuperSet(other Varset) bool {
	return b.BitSet.IsSuperSet(other.(*uibset).BitSet)
}

func (b *uibset) Set(i int) Varset {
	b.BitSet.Set(uint(i))
	return b
}

func (b *uibset) Count() int {
	return int(b.BitSet.Count())
}

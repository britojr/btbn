package varset

import "github.com/willf/bitset"

type uibset struct {
	*bitset.BitSet
}

func newUibset(size int) Varset {
	return &uibset{bitset.New(uint(size))}
}

func (b *uibset) HashString() string {
	return b.DumpAsBits()
}

func (b *uibset) SetHashString(s string) {
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

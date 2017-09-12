package varset

import "testing"

func TestString(t *testing.T) {
	cases := []struct {
		size int
		vars []int
	}{
		{3, []int{}},
		{3, []int{1}},
		{3, []int{0, 1}},
		{3, []int{0, 1, 2}},
		{3, []int{1, 2}},
		{3, []int{0, 2}},
		{527, []int{0, 2, 100, 312, 512}},
	}
	for _, tt := range cases {
		b1 := New(tt.size)
		for _, v := range tt.vars {
			b1.Set(v)
		}
		b1str := b1.DumpHashString()
		b2 := New(tt.size)
		b2.LoadHashString(b1str)
		if !b1.Equal(b2) {
			t.Errorf("wrong mapping of string to varset \n(%v)!=\n(%v)", b1str, b2.DumpHashString())
		}
	}
}

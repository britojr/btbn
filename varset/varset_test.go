package varset

import (
	"reflect"
	"sort"
	"testing"
)

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

func TestInts(t *testing.T) {
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
		b := New(tt.size).SetInts(tt.vars)
		got := b.DumpAsInts()
		sort.Ints(got)
		sort.Ints(tt.vars)
		if !reflect.DeepEqual(tt.vars, got) {
			t.Errorf("wrong vars \n(%v)!=\n(%v)", tt.vars, got)
		}
	}
}

func TestEqual(t *testing.T) {
	cases := []struct {
		size int
		vars []int
		i    int
	}{
		{3, []int{1}, 1},
		{3, []int{0, 1, 2}, 2},
		{527, []int{0, 2, 100, 312, 512}, 100},
	}
	for _, tt := range cases {
		b := New(tt.size).SetInts(tt.vars)
		b2 := b.Clone()
		if !b.Equal(b2) {
			t.Errorf("should be equal (%v)!=(%v)", b, b2)
		}
		b.Clear(tt.i)
		if b.Equal(b2) {
			t.Errorf("should be different (%v)!=(%v)", b, b2)
		}
	}
}

func TestIsSuperSet(t *testing.T) {
	cases := []struct {
		size   int
		supset []int
		set    []int
		is     bool
	}{
		{3, []int{1}, []int{1}, true},
		{3, []int{1}, []int{0}, false},
		{3, []int{0, 1, 2}, []int{2}, true},
		{3, []int{0, 1, 2}, []int{2, 3}, false},
		{527, []int{0, 2, 100, 312, 512}, []int{312, 100, 0, 2}, true},
		{527, []int{0, 2, 100, 312, 512}, []int{312, 100, 0, 1}, false},
	}
	for _, tt := range cases {
		b := New(tt.size).SetInts(tt.supset)
		b2 := New(tt.size).SetInts(tt.set)
		if tt.is != b.IsSuperSet(b2) {
			t.Errorf("wrong (%v) supset (%v) should be (%v)", tt.supset, tt.set, tt.is)
		}
	}
}

package ktree

import (
	"reflect"
	"sort"
	"testing"

	"github.com/britojr/tcc/characteristic"
)

func TestDecodeCharTree(t *testing.T) {
	cases := []struct {
		n, k     int
		iphi     []int
		chartree characteristic.Tree
		cliques  [][]int
		children [][]int
		varin    []int
		varout   []int
	}{
		{
			n: 11, k: 3,
			iphi: []int{0, 10, 9, 3, 4, 5, 6, 7, 1, 2, 8},
			chartree: characteristic.Tree{
				P: []int{-1, 5, 0, 0, 2, 8, 8, 1, 0},
				L: []int{-1, 2, -1, -1, 0, 2, 1, 2, -1},
			},
			cliques: [][]int{
				{1, 2, 8},
				{0, 1, 4, 7},
				{1, 2, 8, 10},
				{1, 2, 8, 9},
				{2, 3, 8, 10},
				{1, 2, 4, 7},
				{1, 5, 7, 8},
				{0, 4, 6, 7},
				{1, 2, 7, 8},
			},
			children: [][]int{
				{2, 3, 8},
				{7},
				{4},
				[]int(nil),
				[]int(nil),
				{1},
				[]int(nil),
				[]int(nil),
				{5, 6},
			},
			varin:  []int{-1, 0, 10, 9, 3, 4, 5, 6, 7},
			varout: []int{-1, 2, -1, -1, 1, 8, 2, 1, -1},
		},
	}
	for _, tt := range cases {
		children, cliques, varin, varout := decodeCharTree(&tt.chartree, tt.iphi, tt.n, tt.k)
		for i := range tt.children {
			got := append([]int(nil), children[i]...)
			sort.Ints(got)
			if !reflect.DeepEqual(tt.children[i], got) {
				t.Errorf("children[%v] (%v)!=(%v)", i, tt.children[i], got)
			}
		}
		for i := range tt.cliques {
			got := append([]int(nil), cliques[i]...)
			sort.Ints(got)
			if !reflect.DeepEqual(tt.cliques[i], got) {
				t.Errorf("cliques[%v] (%v)!=(%v)", i, tt.cliques[i], got)
			}
		}

		if !reflect.DeepEqual(tt.varin, varin) {
			t.Errorf("varin (%v)!=(%v)", tt.varin, varin)
		}
		if !reflect.DeepEqual(tt.varout, varout) {
			t.Errorf("varout (%v)!=(%v)", tt.varout, varout)
		}
	}
}

func TestNewFromDecodedCharTree(t *testing.T) {
	cases := []struct {
		children   [][]int
		cliques    [][]int
		varin      []int
		varout     []int
		answerlist [][]int
	}{
		{
			children: [][]int{
				{2, 3, 8},
				{7},
				{4},
				[]int(nil),
				[]int(nil),
				{1},
				[]int(nil),
				[]int(nil),
				{5, 6},
			},
			cliques: [][]int{
				{1, 2, 8},
				{0, 1, 4, 7},
				{1, 2, 8, 10},
				{1, 2, 8, 9},
				{2, 3, 8, 10},
				{1, 2, 4, 7},
				{1, 5, 7, 8},
				{0, 4, 6, 7},
				{1, 2, 7, 8},
			},
			varin:  []int{-1, 0, 10, 9, 3, 4, 5, 6, 7},
			varout: []int{-1, 2, -1, -1, 1, 8, 2, 1, -1},
			answerlist: [][]int{
				{1, 2, 8, 10, -1, -1},
				{2, 3, 8, 10, 3, 1},
				{1, 2, 8, 9, 9, 10},
				{1, 2, 7, 8, 7, 10},
				{1, 2, 4, 7, 4, 8},
				{1, 5, 7, 8, 5, 2},
				{0, 1, 4, 7, 0, 2},
				{0, 4, 6, 7, 6, 1},
			},
		},
	}
	for _, tt := range cases {
		tk := newFromDecodedCharTree(tt.children, tt.cliques, tt.varin, tt.varout)
		got := getVariablesList(tk)
		for i := range tt.answerlist {
			if !reflect.DeepEqual(tt.answerlist[i], got[i]) {
				t.Errorf("(%v)!=(%v)", tt.answerlist[i], got[i])
			}
		}
	}
}

func getVariablesList(tk *Ktree) (m [][]int) {
	queue, r := []*Ktree{tk}, &Ktree{}
	for len(queue) > 0 {
		r, queue = queue[0], queue[1:]
		line := append([]int(nil), r.Variables()...)
		sort.Ints(line)
		line = append(line, r.VarIn())
		line = append(line, r.VarOut())
		m = append(m, line)
		for _, ch := range r.Children() {
			queue = append(queue, ch)
		}
	}
	return
}

func TestUniformSample(t *testing.T) {
	cases := []struct {
		n, k int
	}{
		{11, 3},
	}
	for _, tt := range cases {
		tk := UniformSample(tt.n, tt.k)
		got := getVariablesList(tk)
		if tt.k+3 != len(got[0]) {
			t.Errorf("wrong k (%v)!=(%v)", tt.k+3, len(got[0]))
		}
		if tt.n-tt.k != len(got) {
			t.Errorf("wrong number of nodes (%v)!=(%v)", tt.n-tt.k, len(got))
		}
	}
}

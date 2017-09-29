package ktree

import (
	"reflect"
	"sort"
	"testing"

	"github.com/britojr/tcc/characteristic"
	"github.com/britojr/tcc/codec"
	"github.com/britojr/tcc/dandelion"
	"github.com/britojr/utl/ints"
)

func TestNew(t *testing.T) {
	adj1 := [][]int{
		{1, 2, 3, 4},
		{2, 4, 5, 6},
		{3},
		[]int(nil),
		{5, 6},
		[]int(nil),
		[]int(nil),
	}
	clqs1 := [][]int{
		{0, 1, 2},
		{0, 2, 3},
		{0, 1, 4},
		{1, 4, 5},
		{1, 4, 6},
	}
	ch1 := New([]int{0, 1, 4}, 4, 2)
	ch1.AddChild(New([]int{1, 4, 5}, 5, 0))
	ch1.AddChild(New([]int{1, 4, 6}, 6, 0))
	tk1 := New([]int{0, 1, 2}, -1, -1)
	tk1.AddChild(New([]int{0, 2, 3}, 3, 1))
	tk1.AddChild(ch1)
	got := helperGetAdjList(tk1, len(adj1))
	for i := range adj1 {
		if !reflect.DeepEqual(adj1[i], got[i]) {
			t.Errorf("wrong adj[%v]: (%v)!=(%v)", i, adj1[i], got[i])
		}
	}
	gotclqs := tk1.AllCliques()
	for i := range clqs1 {
		if !reflect.DeepEqual(clqs1[i], gotclqs[i]) {
			t.Errorf("wrong clq[%v]: (%v)!=(%v)", i, clqs1[i], gotclqs[i])
		}
	}
}

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
		got := helperGetVarsList(tk)
		for i := range tt.answerlist {
			if !reflect.DeepEqual(tt.answerlist[i], got[i]) {
				t.Errorf("(%v)!=(%v)", tt.answerlist[i], got[i])
			}
		}
	}
}

func TestUniformSample(t *testing.T) {
	cases := []struct {
		n, k int
	}{
		{11, 3},
	}
	for _, tt := range cases {
		tk := UniformSample(tt.n, tt.k)
		got := helperGetVarsList(tk)
		if tt.k+1 != len(got[0])-2 {
			t.Errorf("wrong clique size k+1 (%v)!=(%v)", tt.k+1, len(got[0])-2)
		}
		if tt.n-tt.k != len(got) {
			t.Errorf("wrong number of nodes (%v)!=(%v)", tt.n-tt.k, len(got))
		}
	}
}

func TestFromCode(t *testing.T) {
	code1 := []*codec.Code{{
		Q: []int{1, 4},
		S: &dandelion.DandelionCode{
			P: []int{0, 1, 3},
			L: []int{-1, 1, 1},
		}}, {
		Q: []int{1, 2, 6},
		S: &dandelion.DandelionCode{
			P: []int{1, 1, 1, 4},
			L: []int{2, 2, 1, 1},
		}}}
	adj1 := [][][]int{{
		{1, 2, 3, 4},
		{2, 4, 5, 6},
		{3},
		[]int(nil),
		{5, 6},
		[]int(nil),
		[]int(nil),
	}, {
		{1, 2, 3, 4, 5, 6, 7},
		{2, 3, 4, 6, 7, 8},
		{3, 5, 6, 7, 8},
		[]int{5},
		[]int{6},
		[]int(nil),
		[]int{8},
		[]int(nil),
		[]int(nil),
	}}
	cases := []struct {
		n, k int
		code *codec.Code
		adj  [][]int
	}{
		{7, 2, code1[0], adj1[0]},
		{7, 2, code1[0], adj1[0]},
		{9, 3, code1[1], adj1[1]},
	}
	for _, tt := range cases {
		tk := FromCode(tt.code)
		got := helperGetAdjList(tk, tt.n)
		for i := range tt.adj {
			if !reflect.DeepEqual(tt.adj[i], got[i]) {
				t.Errorf("wrong adj[%v]: (%v)!=(%v)", i, tt.adj[i], got[i])
			}
		}
		m := helperGetVarsList(tk)
		if tt.k+1 != len(m[0])-2 {
			t.Errorf("wrong clique size k+1 (%v)!=(%v)", tt.k+1, len(m[0])-2)
		}
		if tt.n-tt.k != len(m) {
			t.Errorf("wrong number of nodes (%v)!=(%v)", tt.n-tt.k, len(m))
		}
	}
}

func TestVarInOut(t *testing.T) {
	cases := []struct {
		n, k int
	}{{7, 2}, {9, 3}, {18, 5}}
	for _, tt := range cases {
		tk := UniformSample(tt.n, tt.k)
		queue := []*Ktree{tk}
		for len(queue) > 0 {
			pa := queue[0]
			queue = queue[1:]
			for _, ch := range pa.Children() {
				vOut := ints.Difference(pa.Variables(), ch.Variables())
				if !reflect.DeepEqual(vOut, []int{ch.VarOut()}) {
					t.Errorf("wrong varOut: (%v)!=(%v)", vOut, ch.VarOut())
				}
				vIn := ints.Difference(ch.Variables(), pa.Variables())
				if !reflect.DeepEqual(vIn, []int{ch.VarIn()}) {
					t.Errorf("wrong varIn: (%v)!=(%v)", vIn, ch.VarIn())
				}
			}
			queue = append(queue, pa.Children()...)
		}
	}
}

// returns clique list with varin varout at the end of each clique
func helperGetVarsList(tk *Ktree) (m [][]int) {
	queue, r := []*Ktree{tk}, &Ktree{}
	for len(queue) > 0 {
		r = queue[0]
		queue = queue[1:]
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

// returns the adjacency matrix of a ktree
func helperGetAdjList(tk *Ktree, n int) [][]int {
	adj := make([][]int, n)
	for i, u := range tk.Variables() {
		for _, v := range tk.Variables()[i+1:] {
			if u < v {
				adj[u] = append(adj[u], v)
			} else {
				adj[v] = append(adj[v], u)
			}
		}
	}
	queue := append([]*Ktree(nil), tk.Children()...)
	for len(queue) > 0 {
		r := queue[0]
		queue = queue[1:]
		u := r.VarIn()
		for _, v := range r.Variables() {
			if v == u {
				continue
			}
			if u < v {
				adj[u] = append(adj[u], v)
			} else {
				adj[v] = append(adj[v], u)
			}
		}
		queue = append(queue, r.Children()...)
	}
	for i := range adj {
		sort.Ints(adj[i])
	}
	return adj
}

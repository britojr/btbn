// Package ktree implements a k-tree structure
// that is, a clique-tree where each node has length k
// and differs from its neighbour in exactly one variable
package ktree

import (
	"fmt"
	"sort"

	"github.com/britojr/tcc/characteristic"
	"github.com/britojr/tcc/codec"
	"github.com/britojr/tcc/generator"
	"github.com/britojr/utl/errchk"
)

// Ktree defines a ktree structure
type Ktree struct {
	clique   []int    // list of variables in a clique
	vIn      int      // variable included relative to the parent node
	vOut     int      // variable removed relative to the parent node
	parent   *Ktree   // pointer to parent
	children []*Ktree // pointer list to chilren
}

// New creates a new ktree node
func New(clique []int, vIn, vOut int) *Ktree {
	tk := new(Ktree)
	tk.clique = clique
	tk.vIn = vIn
	tk.vOut = vOut
	return tk
}

// Variables returns the variables of this node
func (tk *Ktree) Variables() []int {
	return tk.clique
}

// VarIn returns the variable that was included in this node relative to the parent node
func (tk *Ktree) VarIn() int {
	return tk.vIn
}

// VarOut returns the variable that was removed from this node relative to the parent node
func (tk *Ktree) VarOut() int {
	return tk.vOut
}

// Children returns the children nodes of this node
func (tk *Ktree) Children() []*Ktree {
	return tk.children
}

// AddChild adds a child to current node
func (tk *Ktree) AddChild(ch *Ktree) {
	tk.children = append(tk.children, ch)
}

// Parent returns pointer to parent node
func (tk *Ktree) Parent() *Ktree {
	return tk.parent
}

// SetParent updates parent node
func (tk *Ktree) SetParent(pa *Ktree) {
	tk.parent = pa
}

// UniformSample uniformly samples a ktree
func UniformSample(n, k int) *Ktree {
	C, err := generator.RandomCode(n, k)
	errchk.Check(err, "")
	return FromCode(C)
}

// FromCode creates a ktree from a given code
func FromCode(C *codec.Code) *Ktree {
	// Decode a characteristic tree
	T, iphi, err := codec.DecodeCharTree(C)
	errchk.Check(err, "")
	k := len(C.Q)
	n := k + len(T.P) - 1
	return newFromDecodedCharTree(decodeCharTree(T, iphi, n, k))
}

func (tk *Ktree) String() string {
	cl := tk.AllCliques()
	s := fmt.Sprintf("{size: %v cl:{%v}}", len(cl), cl)
	return s
}

// AllCliques returns a list of all cliques that descend from this node
func (tk *Ktree) AllCliques() (cl [][]int) {
	queue := []*Ktree{tk}
	for len(queue) > 0 {
		r := queue[0]
		queue = queue[1:]
		cl = append(cl, r.Variables())
		queue = append(queue, r.Children()...)
	}
	return
}

func decodeCharTree(T *characteristic.Tree, iphi []int, n, k int) (
	[][]int, [][]int, []int, []int,
) {
	// Create children matrix
	children := make([][]int, len(T.P))
	for i := 0; i < len(T.P); i++ {
		if T.P[i] != -1 {
			children[T.P[i]] = append(children[T.P[i]], i)
		}
	}
	// decode clique list
	K, varout := decodeCliqueList(T, children, n, k)
	// relable clique list
	cliques, varin, varout := decodeRelabel(K, varout, iphi)
	return children, cliques, varin, varout
}

func decodeCliqueList(T *characteristic.Tree, children [][]int, n, k int) ([][]int, []int) {
	// Initialize auxiliar (not relabled) clique matrix with the last k variables
	varout := make([]int, len(children))
	K := make([][]int, len(T.P))
	K[0] = make([]int, k)
	for i := 0; i < k; i++ {
		K[0][i] = n - (k - i) + 1
	}
	// Visit T in BFS order, starting with the children of the root
	queue := make([]int, 0, n)
	queue = append(queue, children[0]...)
	for len(queue) > 0 {
		v := queue[0]
		queue = queue[1:]
		// update unlabled clique K[pa[v]]
		for i := 0; i < len(K[T.P[v]]); i++ {
			if i != T.L[v] {
				K[v] = append(K[v], K[T.P[v]][i])
			} else {
				varout[v] = K[T.P[v]][i]
			}
		}
		if T.P[v] != 0 {
			K[v] = append(K[v], T.P[v])
			sort.Ints(K[v])
		}
		// enqueue the children of v
		for i := 0; i < len(children[v]); i++ {
			queue = append(queue, children[v][i])
		}
	}
	return K, varout
}

func decodeRelabel(K [][]int, varout []int, iphi []int) ([][]int, []int, []int) {
	// create relabled cliques list
	cliques := make([][]int, len(K))
	varin := make([]int, len(K))
	varin[0] = -1
	for i := range K {
		cliques[i] = make([]int, len(K[i]))
		for j := range K[i] {
			cliques[i][j] = iphi[K[i][j]-1]
		}
		if i > 0 {
			// append relabled varin
			cliques[i] = append(cliques[i], iphi[i-1])
			varin[i] = iphi[i-1]
		}
	}
	// relable varout
	for i, v := range varout {
		if v > 0 {
			varout[i] = iphi[v-1]
		} else {
			varout[i] = -1
		}
	}
	return cliques, varin, varout
}

func newFromDecodedCharTree(children, cliques [][]int, varin []int, varout []int) *Ktree {
	tk := new(Ktree)
	createNodes(tk, 0, children, cliques, varin, varout)
	tk = removeRoot(tk)
	return tk
}

func createNodes(
	tk *Ktree, i int,
	children, cliques [][]int, varin []int, varout []int,
) {
	tk.clique = cliques[i]
	tk.vIn = varin[i]
	tk.vOut = varout[i]
	for _, v := range children[i] {
		tk.children = append(tk.children, new(Ktree))
		createNodes(tk.children[len(tk.children)-1], v, children, cliques, varin, varout)
	}
}

func removeRoot(tk *Ktree) *Ktree {
	firstChild := tk.children[0]
	vOut := firstChild.vIn
	for _, ch := range tk.children[1:] {
		ch.vOut = vOut
	}
	firstChild.vIn, firstChild.vOut = -1, -1
	firstChild.children = append(firstChild.children, tk.children[1:]...)
	return firstChild
}

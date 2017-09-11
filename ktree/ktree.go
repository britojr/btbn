package ktree

import (
	"sort"

	"github.com/britojr/tcc/characteristic"
	"github.com/britojr/tcc/generator"
	"github.com/britojr/utl/errchk"
)

// Ktree defines a ktree structure
type Ktree struct {
	clique   []int    // list of variables in a clique
	vIn      int      // variable included relative to the parent node
	vOut     int      // variable removed relative to the parent node
	children []*Ktree // pointer list to chilren
}

// Variables returns the variables of this node
func (t *Ktree) Variables() []int {
	panic("not implemented")
}

// Children returns the children of this node
func (t *Ktree) Children() []*Ktree {
	panic("not implemented")
}

// VarIn returns the variable that was included in this node relative to the parent node
func (t *Ktree) VarIn() int {
	panic("not implemented")
}

// VarOut returns the variable that was removed from this node relative to the parent node
func (t *Ktree) VarOut() int {
	panic("not implemented")
}

// UniformSample uniformly samples a ktree
func UniformSample(n, k int) *Ktree {
	T, iphi, err := generator.RandomCharTree(n, k)
	errchk.Check(err, "")
	children, cliques, varin, varout := decodeCharTree(T, iphi, n, k)
	tk := newFromDecodedCharTree(children, cliques, varin, varout)
	return tk
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
	cliques, varin, varout := decodeRelable(K, varout, iphi)
	return children, cliques, varin, varout
}

func decodeCliqueList(T *characteristic.Tree, children [][]int, n, k int) ([][]int, []int) {
	// Initialize auxiliar (not relabled) clique matrix with the last k variables
	varout := make([]int, len(children))
	varout[0] = -1
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

func decodeRelable(K [][]int, varout []int, iphi []int) ([][]int, []int, []int) {
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
	for i := range varout {
		varout[i] = iphi[varout[i]-1]
	}
	return cliques, varin, varout
}

func newFromDecodedCharTree(children, cliques [][]int, varin []int, varout []int) *Ktree {
	r := new(Ktree)
	createNodes(r, 0, children, cliques, varin, varout)
	r = removeRoot(r)
	return r
}

func createNodes(
	r *Ktree, i int,
	children, cliques [][]int, varin []int, varout []int,
) {
	r.clique = cliques[i]
	r.vIn = varin[i]
	r.vOut = varout[i]
	for _, v := range children[i] {
		r.children = append(r.children, new(Ktree))
		createNodes(r.children[len(r.children)-1], v, children, cliques, varin, varout)
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

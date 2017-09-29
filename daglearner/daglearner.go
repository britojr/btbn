package daglearner

import (
	"math/rand"
	"time"

	"github.com/britojr/btbn/bnstruct"
	"github.com/britojr/btbn/ktree"
	"github.com/britojr/btbn/scr"
	"github.com/britojr/btbn/varset"
)

var seed = func() int64 {
	return time.Now().UnixNano()
}

// Approximated learns a dag approximatedly from a ktree
func Approximated(tk *ktree.Ktree, ranker scr.Ranker) (bn *bnstruct.BNStruct) {
	// Initialize the local scores for empty list of parents
	bn = bnstruct.New(ranker.Size())
	parents := varset.New(ranker.Size())
	for x := 0; x < ranker.Size(); x++ {
		bn.SetParents(x, parents, ranker.ScoreOf(x, parents))
	}

	// Sample a partial order from the ktree
	pOrders := samplePartialOrder(tk)

	// Find the parents sets that maximize the score and respect the partial order
	for _, pOrd := range pOrders {
		setParentsFromOrder(pOrd, ranker, bn)
	}
	return bn
}

func setParentsFromOrder(order partialOrder, ranker scr.Ranker, bn *bnstruct.BNStruct) {
	restric := varset.New(ranker.Size())
	for _, v := range order.vars[:order.ini] {
		restric.Set(v)
	}
	for _, v := range order.vars[order.ini:] {
		newParents, newScore := ranker.BestIn(v, restric)
		if newScore > bn.LocalScore(v) {
			bn.SetParents(v, newParents, newScore)
		}
		restric.Set(v)
	}
}

type partialOrder struct {
	vars []int
	ini  int
}

// samplePartialOrder samples a partial order from a ktree
func samplePartialOrder(tk *ktree.Ktree) []partialOrder {
	r := rand.New(rand.NewSource(seed()))
	// start partial order with a shuffle of the root node
	po := []partialOrder{
		partialOrder{shuffle(tk.Variables(), r), 0},
	}
	// for each children node, theres is one variable replaced relative to its parent
	// sample a position to insert the respective variable, respecting the prevoious orders
	sampleChildrenOrder(tk, po[0].vars, r, &po)
	return po
}

func sampleChildrenOrder(tk *ktree.Ktree, paOrder []int, r *rand.Rand, po *[]partialOrder) {
	for _, ch := range tk.Children() {
		chOrder := getChildOrder(paOrder, ch.VarIn(), ch.VarOut(), r)
		(*po) = append((*po), chOrder)
		sampleChildrenOrder(ch, chOrder.vars, r, po)
	}
}

func getChildOrder(paOrder []int, vIn, vOut int, r *rand.Rand) partialOrder {
	pos := r.Intn(len(paOrder))
	chOrder := make([]int, len(paOrder))
	// Must use the 'for' without increment
	// because the 'continue' should not increment i
	for i, j := 0, 0; i < len(chOrder); {
		switch {
		case j < len(paOrder) && paOrder[j] == vOut:
			j++
		case i == pos:
			chOrder[i] = vIn
			i++
		default:
			chOrder[i] = paOrder[j]
			i, j = i+1, j+1
		}
	}
	return partialOrder{chOrder, pos}
}

func shuffle(xs []int, r *rand.Rand) []int {
	perm := r.Perm(len(xs))
	shuf := make([]int, len(xs))
	for i := range perm {
		shuf[i] = xs[perm[i]]
	}
	return shuf
}

// Exact learns an optimal dag from a ktree
// TODO: need to replace this for an actual call to an exact method
func Exact(tk *ktree.Ktree, ranker scr.Ranker) (bn *bnstruct.BNStruct) {
	bn = Approximated(tk, ranker)
	for i := 0; i < 50; i++ {
		currBn := Approximated(tk, ranker)
		if currBn.Better(bn) {
			bn = currBn
		}
	}
	return
}

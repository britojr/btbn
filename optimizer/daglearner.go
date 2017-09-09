package optimizer

import (
	"math/rand"
	"time"

	"github.com/britojr/btbn/ctree"
	"github.com/britojr/btbn/score"
	"github.com/britojr/btbn/varset"
)

var seed = func() int64 {
	return time.Now().UnixNano()
}

// DAGapproximatedLearning learns a dag approximatedly from a  ktree
func DAGapproximatedLearning(c *ctree.Ctree, rankers []score.Ranker) (bn *BNStructure) {
	// Initialize the local scores for empty list of parents
	bn = NewBNStructure()
	parents := varset.New(len(rankers))
	for x := 0; x < len(rankers); x++ {
		bn.SetParents(x, parents, rankers[x].ScoreOf(parents))
	}

	// Sample a partial order from the ktree
	pOrders := samplePartialOrder(c, seed())

	// Find the parents sets that maximize the score and respect the partial order
	for _, pOrd := range pOrders {
		setParentsFromOrder(pOrd, rankers, bn)
	}
	return bn
}

func setParentsFromOrder(order partialOrder, rankers []score.Ranker, bn *BNStructure) {
	restric := varset.New(len(rankers))
	for _, v := range order.vars[:order.ini] {
		restric.Set(v)
	}
	for _, v := range order.vars[order.ini:] {
		newParents, newScore := rankers[v].BestIn(restric)
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
func samplePartialOrder(c *ctree.Ctree, seed int64) []partialOrder {
	r := rand.New(rand.NewSource(seed))
	// start partial order with a shuffle of the root node
	po := []partialOrder{
		partialOrder{shuffle(c.Variables(), r), 0},
	}
	// for each children node, theres is one variable replaced relative to its parent
	// sample a position to insert the respective variable, respecting the prevoious orders
	sampleChildrenOrder(c, po[0].vars, r, &po)
	return po
}

func sampleChildrenOrder(c *ctree.Ctree, paOrder []int, r *rand.Rand, po *[]partialOrder) {
	for _, ch := range c.Children() {
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

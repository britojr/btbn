package daglearner

import (
	"math/rand"

	"github.com/britojr/btbn/ctree"
	"github.com/britojr/btbn/optimizer"
	"github.com/britojr/btbn/score"
	"github.com/britojr/btbn/varset"
)

// Approximated learns a dag approximatedly from a  ktree
func Approximated(ct *ctree.Ctree, rankers []score.Ranker) (bn *optimizer.BNStructure) {

	// Initialize the local scores for empty list of parents
	bn = optimizer.NewBNStructure()
	parents := varset.New(len(rankers))
	for x := 0; x < len(rankers); x++ {
		bn.SetParents(x, parents, rankers[x].ScoreOf(parents))
	}

	restric := varset.New(len(rankers))
	variables := ct.Variables()
	perm := rand.Perm(len(variables))
	restric.Set(variables[perm[0]])
	varOrder := make([]int, len(variables)-1)
	// Select parents for the second variable forward
	for i := 1; i < len(perm); i++ {
		varOrder[i] = variables[perm[i]]
	}
	findBestParents(varOrder, restric, rankers, bn)
	orderChildren(varOrder, ct, rankers, bn)
	// Visit tree starting with the root's children
	// queue := make([]*ctree.Ctree, 0)
	// for _, ch := range ct.Children() {
	// 	queue = append(queue, ch)
	// }
	// queue = append(queue, ct.Children()...)
	// for len(queue) > 0 {
	// node := queue[0]
	// restric := varset.New(len(rankers))
	// variables := node.SepSet()
	// }

	return bn
}

func orderChildren(
	pVarOrder []int,
	node *ctree.Ctree,
	rankers []score.Ranker,
	bn *optimizer.BNStructure,
) {
	for _, ch := range node.Children() {
		restric := varset.New(len(rankers))
		pos := rand.Intn(len(pVarOrder))
		cVarOrder := make([]int, len(pVarOrder))
		i, j := 0, 0
		// Must use the 'for' without increment because the 'continue' should
		// not increment i
		for i < len(cVarOrder) {
			if j < len(pVarOrder) && pVarOrder[j] == ch.VarOut() {
				j++
				continue
			}
			if i == pos {
				cVarOrder[i] = ch.VarIn()
				findBestParents([]int{cVarOrder[i]}, restric, rankers, bn)
			} else {
				cVarOrder[i] = pVarOrder[j]
				j++
			}
			if i < pos {
				restric.Set(cVarOrder[i])
			}
			i++
		}
		findBestParents(cVarOrder[pos+1:], restric, rankers, bn)
		orderChildren(cVarOrder, ch, rankers, bn)
	}
}

func findBestParents(
	varOrder []int,
	restric varset.Varset,
	rankers []score.Ranker,
	bn *optimizer.BNStructure,
) {
	for _, v := range varOrder {
		newParents, newScore := rankers[v].BestIn(restric)
		if newScore > bn.LocalScore(v) {
			bn.SetParents(v, newParents, newScore)
		}
		restric.Set(v)
	}
}

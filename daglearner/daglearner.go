package daglearner

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"math/rand"
	"path/filepath"
	"strings"
	"time"

	"github.com/britojr/btbn/bnstruct"
	"github.com/britojr/btbn/ktree"
	"github.com/britojr/btbn/scr"
	"github.com/britojr/btbn/varset"
	"github.com/britojr/utl/cmdsh"
	"github.com/britojr/utl/conv"
	"github.com/britojr/utl/errchk"
	"github.com/britojr/utl/ioutl"
)

var seed = func() int64 {
	return time.Now().UnixNano()
}

// Approximated learns a dag approximatedly from a ktree
func Approximated(tk *ktree.Ktree, ranker scr.Ranker) (bn *bnstruct.BNStruct) {
	bn = buildEmptyBN(ranker)

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
		{shuffle(tk.Variables(), r), 0},
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

func buildEmptyBN(ranker scr.Ranker) *bnstruct.BNStruct {
	// Initialize the local scores for empty list of parents
	bn := bnstruct.New(ranker.Size())
	parents := varset.New(ranker.Size())
	for x := 0; x < ranker.Size(); x++ {
		bn.SetParents(x, parents, ranker.ScoreOf(x, parents))
	}
	return bn
}

// Exact learns an optimal dag from a ktree
func Exact(tk *ktree.Ktree, ranker scr.Ranker) (bn *bnstruct.BNStruct) {
	f, err := ioutil.TempFile("", "gob-")
	errchk.Check(err, "")
	fname := filepath.Base(f.Name())
	f.Close()
	ranker.SaveSubSet(fname, tk.Variables())
	_, err = cmdsh.Exec(fmt.Sprintf("gobnilp -f=pss %s", fname), 0)
	errchk.Check(err, "")
	solFile := strings.TrimSuffix(fname, filepath.Ext(fname)) + ".solution"
	paLst, paScr := parseParentMat(solFile)

	bn = buildEmptyBN(ranker)
	for v, pa := range paLst {
		paset := varset.New(ranker.Size())
		paset.SetInts(pa)
		bn.SetParents(v, paset, paScr[v])
	}
	return
}

func parseParentMat(fname string) (map[int][]int, map[int]float64) {
	paLst := make(map[int][]int)
	paScr := make(map[int]float64)
	r := ioutl.OpenFile(fname)
	defer r.Close()
	scanner := bufio.NewScanner(r)
	paSep := "<-"
	for scanner.Scan() {
		text := strings.Replace(scanner.Text(), ":", paSep, 1)
		text = strings.Replace(text, ",", " ", -1)
		line := strings.SplitN(text, paSep, 2)
		if len(line) < 2 {
			continue
		}
		vID := conv.Atoi(line[0])
		fields := strings.Fields(line[1])
		paLst[vID] = conv.Satoi(fields[:len(fields)-1])
		paScr[vID] = conv.Atof(fields[len(fields)-1])
	}
	return paLst, paScr
}

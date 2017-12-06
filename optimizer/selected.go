package optimizer

import (
	"log"
	"math"
	"math/rand"
	"sort"

	"github.com/britojr/btbn/bnstruct"
	"github.com/britojr/btbn/daglearner"
	"github.com/britojr/btbn/ktree"
	"github.com/britojr/btbn/scr"
	"github.com/britojr/btbn/varset"
	"github.com/britojr/tcc/codec"
	"github.com/britojr/tcc/generator"
	"github.com/britojr/utl/conv"
	"github.com/britojr/utl/errchk"
	"github.com/britojr/utl/stats"
	"github.com/gonum/floats"
)

// SelectSampleSearch implements the select sampling strategy
type SelectSampleSearch struct {
	*common                   // common variables and methods
	mutInfoFile string        // file with pre-computed mutual info
	prevCodes   []*codec.Code // previously accepted codes
	bestIScr    float64       // currently best IScore
	numTrees    int           // number of ktrees to sample before start learning DAG
	tkList      []*scr.Record // list of accepted ktrees sorted by score
	mutInfo     *scr.MutInfo  // pre-computed mutual information matrix
	kernelZero  float64       // pre-calculated kernel(0)
}

// NewSelectSampleSearch creates a instance of the selected sample stragegy
func NewSelectSampleSearch(scoreRanker scr.Ranker) Optimizer {
	s := &SelectSampleSearch{common: newCommon(scoreRanker)}
	s.kernelZero = stats.GaussianKernel(0.0)
	return s
}

// Search searches for a network structure
func (s *SelectSampleSearch) Search() *bnstruct.BNStruct {
	if len(s.tkList) == 0 {
		s.selectKTrees()
	}
	bn := daglearner.Approximated(s.tkList[0].Data().(*ktree.Ktree), s.scoreRanker)
	s.tkList = s.tkList[1:]
	return bn
}

// SetDefaultParameters sets parameters to default values
func (s *SelectSampleSearch) SetDefaultParameters() {
	s.common.SetDefaultParameters()
	s.numTrees = 1
}

// SetFileParameters sets parameters from input file
func (s *SelectSampleSearch) SetFileParameters(parms map[string]string) {
	s.common.SetFileParameters(parms)
	if numTrees, ok := parms[ParmNumTrees]; ok {
		s.numTrees = conv.Atoi(numTrees)
	}
	if mutInfoFile, ok := parms[ParmMutualInfo]; ok {
		s.mutInfoFile = mutInfoFile
	}
}

// ValidateParameters validates internal parameters
func (s *SelectSampleSearch) ValidateParameters() {
	s.common.ValidateParameters()
	if len(s.mutInfoFile) == 0 {
		log.Panic("Mutual information file missing")
	}
	s.mutInfo = scr.ReadMutInf(s.mutInfoFile)
}

// PrintParameters prints the algorithm's current parameters
func (s *SelectSampleSearch) PrintParameters() {
	s.common.PrintParameters()
	log.Printf("%v: %v\n", ParmNumTrees, s.numTrees)
	log.Printf("%v: '%v'\n", ParmMutualInfo, s.mutInfoFile)
}

// selectKTrees samples and selects a given number of ktrees
func (s *SelectSampleSearch) selectKTrees() {
	r := rand.New(rand.NewSource(seed()))
	s.tkList = make([]*scr.Record, 0, s.numTrees)
	for len(s.tkList) < s.numTrees {
		C, err := generator.RandomCode(s.nv, s.tw)
		errchk.Check(err, "")
		if !s.acceptCode(C, r) {
			continue
		}
		tk := ktree.FromCode(C)
		iscr := s.computeIScore(tk)
		if !s.acceptTree(iscr, r) {
			continue
		}
		s.tkList = append(s.tkList, scr.NewRecord(iscr, tk))
	}
	sort.Slice(s.tkList, func(i, j int) bool {
		return s.tkList[i].Score() > s.tkList[j].Score()
	})
}

func (s *SelectSampleSearch) acceptTree(iscr float64, r *rand.Rand) bool {
	if iscr > s.bestIScr {
		s.bestIScr = iscr
		return true
	}
	return r.Float64() <= (iscr / s.bestIScr)
}

func (s *SelectSampleSearch) computeIScore(tk *ktree.Ktree) float64 {
	var mi float64
	partialScores := make([]float64, s.nv)
	restric := varset.New(s.nv)
	s.computeNodeIScore(tk, restric, partialScores, &mi)
	return mi / math.Abs(floats.Sum(partialScores))
}

func (s *SelectSampleSearch) computeNodeIScore(
	nd *ktree.Ktree, restric varset.Varset, partialScores []float64, mi *float64,
) {
	if nd.VarIn() >= 0 {
		u := nd.VarIn()
		restric.Set(u)
		restric.Clear(nd.VarOut())
		for _, v := range nd.Variables() {
			if v != u {
				*mi += s.mutInfo.Get(u, v)
			}
		}
	} else {
		// if got here, it is the root node, then initialize the restriction set
		for i, u := range nd.Variables() {
			restric.Set(u)
			for _, v := range nd.Variables()[i+1:] {
				*mi += s.mutInfo.Get(u, v)
			}
		}
	}
	for _, v := range nd.Variables() {
		_, newScore := s.scoreRanker.BestIn(v, restric)
		if partialScores[v] == 0 || newScore > partialScores[v] {
			partialScores[v] = newScore
		}
	}
	for _, ch := range nd.Children() {
		s.computeNodeIScore(ch, restric, partialScores, mi)
	}
	if nd.VarOut() >= 0 {
		restric.Set(nd.VarOut())
		restric.Clear(nd.VarIn())
	}
}

// acceptCode stochastically accepts a Dandelion code
func (s *SelectSampleSearch) acceptCode(C *codec.Code, r *rand.Rand) bool {
	// If not the first code, compute the probability of accepting
	if len(s.prevCodes) != 0 {
		if r.Float64() > s.acceptCodeProb(C) {
			return false
		}
	}
	s.prevCodes = append(s.prevCodes, C)
	return true
}

// acceptCodeProb calculates the probability of accepting a code
// based on its distance from the previous ones
func (s *SelectSampleSearch) acceptCodeProb(C *codec.Code) float64 {
	q := float64(0)
	for _, prevCode := range s.prevCodes {
		q += stats.GaussianKernel(CodeDistance(C, prevCode))
	}
	q /= float64(len(s.prevCodes))
	return 1.0 - (q / s.kernelZero)
}

// CodeDistance calculates the distance between two dandelion codes:
//		||C1-C2|| = ||C1.Q - C2.Q||_2 + ||C1.S - C2.S||_2,1
func CodeDistance(C1, C2 *codec.Code) float64 {
	dq := stats.IntsLNormDiff(C1.Q, C2.Q, 2)
	dp := stats.IntsLNormDiff(C1.S.P, C2.S.P, 2)
	dl := stats.IntsLNormDiff(C1.S.L, C2.S.L, 2)
	// return dq + math.Abs(dp-dl)
	return dq + dp + dl
}

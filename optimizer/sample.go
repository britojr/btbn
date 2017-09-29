package optimizer

import (
	"github.com/britojr/btbn/bnstruct"
	"github.com/britojr/btbn/daglearner"
	"github.com/britojr/btbn/ktree"
	"github.com/britojr/btbn/scr"
)

// SampleSearch implements the sampling strategy
type SampleSearch struct {
	*common // common variables and methods
}

// NewSampleSearch creates a instance of the sample stragegy
func NewSampleSearch(scoreRanker scr.Ranker) Optimizer {
	return &SampleSearch{common: newCommon(scoreRanker)}
}

// Search searchs for a network structure
func (s *SampleSearch) Search() *bnstruct.BNStruct {
	tk := ktree.UniformSample(s.nv, s.tw)
	bn := daglearner.Approximated(tk, s.scoreRanker)
	return bn
}

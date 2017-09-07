package optimizer

import (
	"github.com/britojr/btbn/score"
	"github.com/britojr/utl/conv"
)

// SampleSearch implements the sampling strategy
type sampleSearch struct {
	tw int //treewidth
}

// NewSampleSearch creates a instance of the sample stragegy
func newSampleSearch(scoreRankers []score.Ranker, parmFile string) Optimizer {
	s := new(sampleSearch)
	setParameters(s, parmFile)
	return s
}

// Search search network structure
func (s *sampleSearch) Search(numSolutions, timeAvailable int) interface{} {
	return nil
}

// SetDefaultParameters sets the defaults
func (s *sampleSearch) SetDefaultParameters() {
	// set internal variables to defined constants
	s.tw = 3
}

// SetFileParameters sets parameters from input file
func (s *sampleSearch) SetFileParameters(parms map[string]string) {
	if tw, ok := parms["treewidth"]; ok {
		s.tw = conv.Atoi(tw)
	}
}

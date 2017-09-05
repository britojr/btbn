package optimizer

import "github.com/britojr/btbn/score"

// SampleSearch implements the sampling strategy
type SampleSearch struct {
}

// Search search network structure
func (s *SampleSearch) Search(numSolutions, timeAvailable int) interface{} {
	return nil
}

// NewSampleSearch ..
func NewSampleSearch(parmFile string, scoreRankers []score.Ranker) *SampleSearch {
	return nil
}

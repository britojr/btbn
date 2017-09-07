package optimizer

import "github.com/britojr/btbn/score"

// SampleSearch implements the sampling strategy
type SampleSearch struct {
}

// NewSampleSearch creates a instance of the sample stragegy
func NewSampleSearch(scoreRankers []score.Ranker, parmFile string) *SampleSearch {
	return nil
}

// Search search network structure
func (s *SampleSearch) Search(numSolutions, timeAvailable int) interface{} {
	return nil
}

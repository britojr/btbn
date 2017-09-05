package score

// Ranker defines a list of best scores for a given variable
type Ranker interface {
}

// CreateRankers creates array of rankers, one for each variable
func CreateRankers(cache *Cache, maxPa int) []Ranker {
	return nil
}

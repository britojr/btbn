package score

// Cache stores pre-computed score information
type Cache struct {
}

// Read reads a score file into a score cache
func Read(fname string) *Cache {
	return &Cache{}
}

// Nvar returns the number of variables
func (c *Cache) Nvar() int {
	return 0
}

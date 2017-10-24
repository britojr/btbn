package vars

// Index defines a way to iterate through joints of variables
type Index struct {
	current       int
	attrb, stride []int
	vs            VarList
}

// NewIndexFor creates a new index to iterate through indexVars relative to forVars
func NewIndexFor(indexVars, forVars VarList) (ix *Index) {
	ix = new(Index)
	ix.vs = forVars
	ix.attrb = make([]int, len(forVars))
	ix.stride = make([]int, len(forVars))
	j := 0
	s := 1
	for _, v := range indexVars {
		for ; j < len(forVars) && forVars[j].ID() < v.ID(); j++ {
		}
		if j < len(forVars) && forVars[j].ID() == v.ID() {
			ix.stride[j] = s
			j++
		}
		s *= v.NState()
	}
	return
}

// I returns current index
func (ix *Index) I() int {
	return ix.current
}

// Ended if index reached end value
func (ix *Index) Ended() bool {
	return ix.current < 0
}

// Reset set index to begining value
func (ix *Index) Reset() {
	ix.current = 0
	for i := range ix.attrb {
		ix.attrb[i] = 0
	}
}

// Next iterates to next index value
// returns true if it found a valid index
func (ix *Index) Next() bool {
	if ix.current >= 0 {
		for i := range ix.attrb {
			ix.current += ix.stride[i]
			ix.attrb[i]++
			if ix.attrb[i] < ix.vs[i].NState() {
				return true
			}
			ix.current -= ix.stride[i] * ix.vs[i].NState()
			ix.attrb[i] = 0
		}
		ix.current = -1
	}
	return false
}

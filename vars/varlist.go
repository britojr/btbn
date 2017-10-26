package vars

import "sort"

// VarList an ordered set of variables
type VarList []*Var

// NewList creates a new varlist from variables ids and number of states
func NewList(vs, ns []int) (w VarList) {
	for i, v := range vs {
		if i < len(ns) {
			w = append(w, &Var{v, ns[i]})
		} else {
			w = append(w, &Var{v, DefaultNState})
		}
	}
	sort.Slice(w, func(i int, j int) bool {
		return w[i].ID() < w[j].ID()
	})
	return
}

// NStates returns the number of states of the joint of variables
func (vl VarList) NStates() int {
	s := 1
	for _, v := range vl {
		s *= v.NState()
	}
	return s
}

// Copy returns a copy of vl
// note: although it returns a new slice, the content (i.e. variable pointers) is shared
func (vl VarList) Copy() (w VarList) {
	w = make([]*Var, len(vl))
	copy(w, vl)
	return
}

// Equal returns true if vl is equal to other
func (vl VarList) Equal(other VarList) bool {
	if len(vl) != len(other) {
		return false
	}
	for i, v := range other {
		if vl[i].ID() != v.ID() {
			return false
		}
	}
	return true
}

// Diff returns new list with elements in vl and not in other
func (vl VarList) Diff(other VarList) (w VarList) {
	w = make([]*Var, 0, len(vl))
	j := 0
	for _, v := range vl {
		for ; j < len(other) && other[j].ID() < v.ID(); j++ {
		}
		if j < len(other) && other[j].ID() == v.ID() {
			j++
			continue
		}
		w = append(w, v)
	}
	return
}

// Union returns new list merging elements in vl and in other
func (vl VarList) Union(other VarList) (w VarList) {
	w = make([]*Var, 0, len(vl)+len(other))
	j := 0
	for _, v := range vl {
		for ; j < len(other) && other[j].ID() < v.ID(); j++ {
			w = append(w, other[j])
		}
		if j < len(other) && other[j].ID() == v.ID() {
			w = append(w, v)
			j++
			continue
		}
		w = append(w, v)
	}
	for ; j < len(other); j++ {
		w = append(w, other[j])
	}
	return
}

// IntersecID returns new list with elements present in vl and in ids
func (vl VarList) IntersecID(ids ...int) (w VarList) {
	sort.Ints(ids)
	w = make([]*Var, 0, len(vl))
	j := 0
	for _, v := range vl {
		for ; j < len(ids) && ids[j] < v.ID(); j++ {
		}
		if j < len(ids) && ids[j] == v.ID() {
			w = append(w, v)
			j++
			continue
		}
	}
	return
}

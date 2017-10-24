package vars

import "testing"

func TestCopy(t *testing.T) {
	cases := []struct {
		vl VarList
	}{
		{[]*Var{}},
		{[]*Var(nil)},
		{[]*Var{{0, 2}, {3, 2}, {5, 4}}},
	}
	for _, tt := range cases {
		w := tt.vl.Copy()
		if !tt.vl.Equal(w) {
			t.Errorf("not equal %v != %v", tt.vl, w)
		}
		// test if its safe for change
		if len(tt.vl) > 0 {
			w[0] = &Var{9, 2}
		} else {
			w = append(w, &Var{9, 2})
			tt.vl = append(tt.vl, &Var{8, 2})
		}
		if tt.vl.Equal(w) {
			t.Errorf("both are pointing to the same slice %v == %v", &tt.vl, &w)
		}
	}
}

func TestNStates(t *testing.T) {
	cases := []struct {
		vl VarList
		ns int
	}{
		{[]*Var{}, 1},
		{[]*Var(nil), 1},
		{[]*Var{{0, 2}, {3, 2}, {5, 4}}, 16},
	}
	for _, tt := range cases {
		got := tt.vl.NStates()
		if tt.ns != got {
			t.Errorf("wrong number of states %v != %v", tt.ns, got)
		}
	}
}

func TestEqual(t *testing.T) {
	cases := []struct {
		va, vb VarList
		eq     bool
	}{
		{[]*Var{}, []*Var{}, true},
		{[]*Var(nil), []*Var(nil), true},
		{[]*Var{{0, 2}, {3, 2}, {5, 4}}, []*Var{{0, 2}, {3, 2}, {5, 4}}, true},
		{[]*Var{{5, 4}}, []*Var{{5, 3}}, true},
		{[]*Var{{0, 2}}, []*Var{{0, 2}, {1, 2}}, false},
		{[]*Var{{0, 2}, {2, 2}}, []*Var{{0, 2}, {1, 2}}, false},
	}
	for _, tt := range cases {
		got := tt.va.Equal(tt.vb)
		if tt.eq != got {
			t.Errorf("wrong compare result %v != %v", tt.eq, got)
		}
		if tt.vb.Equal(tt.va) != got {
			t.Errorf("equal fucntion should be simetric %v != %v", tt.vb.Equal(tt.va), got)
		}
	}
}

func TestNewList(t *testing.T) {
	cases := []struct {
		vs, ns []int
		res    VarList
	}{
		{[]int{1, 5, 0}, []int{3, 2, 4}, []*Var{{0, 4}, {1, 3}, {5, 2}}},
		{[]int{1, 5, 0}, []int(nil), []*Var{{0, 2}, {1, 2}, {5, 2}}},
	}
	for _, tt := range cases {
		got := NewList(tt.vs, tt.ns)
		if !tt.res.Equal(got) {
			t.Errorf("wrong new list %v != %v", tt.res, got)
		}
	}
}

func TestDiff(t *testing.T) {
	cases := []struct {
		a, b, want []int
	}{
		{[]int(nil), []int(nil), []int(nil)},
		{[]int{2, 3}, []int{}, []int{2, 3}},
		{[]int{}, []int{3, 4}, []int(nil)},
		{[]int{1, 2, 3, 4, 5, 6}, []int{2, 4, 6, 8}, []int{1, 3, 5}},
		{[]int{3, 4, 5, 6}, []int{1, 2, 4, 6, 8}, []int{3, 5}},
	}
	for _, tt := range cases {
		a := NewList(tt.a, nil)
		b := NewList(tt.b, nil)
		got := a.Diff(b)
		for i, v := range tt.want {
			if v != got[i].ID() {
				t.Errorf("(%v) - (%v) = (%v) !=(%v)", a, b, tt.want, got)
			}
		}
	}
}

func TestUnion(t *testing.T) {
	cases := []struct {
		a, b, want []int
	}{
		{[]int(nil), []int(nil), []int(nil)},
		{[]int{1}, []int{1}, []int{1}},
		{[]int{3}, []int{1}, []int{1, 3}},
		{[]int{1}, []int{3}, []int{1, 3}},
		{[]int{2, 3}, []int{}, []int{2, 3}},
		{[]int{}, []int{3, 4}, []int{3, 4}},
		{[]int{6, 4, 2, 8}, []int{8, 9, 3, 1, 2}, []int{1, 2, 3, 4, 6, 8, 9}},
	}
	for _, tt := range cases {
		a := NewList(tt.a, nil)
		b := NewList(tt.b, nil)
		got := a.Union(b)
		for i, v := range tt.want {
			if v != got[i].ID() {
				t.Errorf("(%v) union (%v) = (%v) !=(%v)", a, b, tt.want, got)
			}
		}
	}
}

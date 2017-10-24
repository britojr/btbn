package vars

import "testing"

func TestNewIndexFor(t *testing.T) {
	cases := []struct {
		indexVars, forVars VarList
		seq                []int
	}{
		{[]*Var{{0, 2}, {5, 3}}, []*Var{{0, 2}, {3, 2}, {5, 3}}, []int{0, 1, 0, 1, 2, 3, 2, 3, 4, 5, 4, 5}},
	}
	for _, tt := range cases {
		ix := NewIndexFor(tt.indexVars, tt.forVars)
		for i, v := range tt.seq {
			if ix.Ended() {
				t.Errorf("index ended at %v iterations, should have %v", i, len(tt.seq))
			}
			if v != ix.I() {
				t.Errorf("wrong index i[%v]=%v, should be i[%v]%v", i, ix.I(), i, tt.seq[i])
			}
			ix.Next()
		}
		ix.Reset()
		for i, v := range tt.seq {
			if ix.Ended() {
				t.Errorf("RESET index ended at %v iterations, should have %v", i, len(tt.seq))
			}
			if v != ix.I() {
				t.Errorf("wrong RESET index i[%v]=%v, should be i[%v]%v", i, ix.I(), i, tt.seq[i])
			}
			ix.Next()
		}
	}
}

package factor

import (
	"testing"

	"github.com/britojr/btbn/vars"
	"github.com/gonum/floats"
)

func TestNew(t *testing.T) {
	cases := []struct {
		vs vars.VarList
	}{
		{[]*vars.Var{}},
		{[]*vars.Var{vars.New(0, 2)}},
		{[]*vars.Var{vars.New(1, 4), vars.New(3, 2)}},
	}
	for _, tt := range cases {
		got := New(tt.vs)
		if !tt.vs.Equal(got.vs) {
			t.Errorf("wrong vars %v != %v", tt.vs, got.vs)
		}
		if tt.vs.NStates() != len(got.values) {
			t.Errorf("wrong size of values %v != %v", tt.vs.NStates(), got.values)
		}
		if floats.Sum(got.values) != 1 {
			t.Errorf("factor not normalized, sum == %v", floats.Sum(got.values))
		}
	}
}

// func TestNew(t *testing.T) {
// 	cases:= []struct{
//
// 	}{}
// 	for _, tt:= range cases{
//
// 	}
// }

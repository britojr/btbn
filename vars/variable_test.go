package vars

import "testing"

func TestNew(t *testing.T) {
	cases := []struct {
		id, nstate int
	}{
		{1, 3},
		{0, 2},
		{0, 0},
	}
	for _, tt := range cases {
		v := New(tt.id, tt.nstate)
		if tt.id != v.ID() {
			t.Errorf("wrong id %v != %v", tt.id, v.ID())
		}
		if tt.nstate != v.NState() {
			t.Errorf("wrong nstates %v != %v", tt.nstate, v.NState())
		}
	}
}

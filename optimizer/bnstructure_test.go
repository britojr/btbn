package optimizer

import (
	"math"
	"testing"

	"github.com/britojr/btbn/scr"
	"github.com/britojr/btbn/varset"
)

func TestNewBNStructure(t *testing.T) {
	cases := []struct {
		size int
	}{{1}, {2}, {7}, {55}}
	for _, tt := range cases {
		got := NewBNStructure(tt.size)
		if tt.size != got.Size() {
			t.Errorf("wrong size (%v)!=(%v)", tt.size, got.Size())
		}
	}
}

func TestScore(t *testing.T) {
	cases := []struct {
		size   int
		vars   map[int]*scr.Record
		totScr float64
	}{{
		5, map[int]*scr.Record{
			0: scr.NewRecord(-10, varset.New(0)),
			1: scr.NewRecord(-15, varset.New(0)),
			2: scr.NewRecord(-20, varset.New(0)),
			3: scr.NewRecord(-120, varset.New(0)),
			4: scr.NewRecord(-2, varset.New(0)),
		}, -167,
	}, {
		5, map[int]*scr.Record{
			1: scr.NewRecord(-15, varset.New(0)),
			2: scr.NewRecord(-20, varset.New(0)),
			4: scr.NewRecord(-10, varset.New(0)),
		}, math.Inf(-1),
	}}
	for _, tt := range cases {
		bn := NewBNStructure(tt.size)
		for v, r := range tt.vars {
			bn.SetParents(v, r.Data().(varset.Varset), r.Score())
		}
		for v := 0; v < tt.size; v++ {
			scr := bn.LocalScore(v)
			r, ok := tt.vars[v]
			if ok {
				if r.Score() != scr {
					t.Errorf("wrong local score of (%v): (%v)!=(%v)", v, r.Score(), scr)
				}
			} else {
				if tt.totScr != scr {
					t.Errorf("wrong local score of (%v): (%v)!=(%v)", v, tt.totScr, scr)
				}
			}
		}
		scr := bn.Score()
		if tt.totScr != scr {
			t.Errorf("wrong total score (%v)!=(%v)", tt.totScr, scr)
		}
	}
}

func TestParents(t *testing.T) {
	cases := []struct {
		size int
		vars map[int]*scr.Record
	}{{
		5, map[int]*scr.Record{
			0: scr.NewRecord(-10, varset.New(1).Set(0)),
			1: scr.NewRecord(-15, varset.New(0)),
			2: scr.NewRecord(-20, varset.New(0)),
			3: scr.NewRecord(-120, varset.New(0)),
			4: scr.NewRecord(-2, varset.New(0)),
		},
	}, {
		5, map[int]*scr.Record{
			1: scr.NewRecord(-15, varset.New(0)),
			2: scr.NewRecord(-20, varset.New(0)),
			4: scr.NewRecord(-10, varset.New(0)),
		},
	}}
	for _, tt := range cases {
		bn := NewBNStructure(tt.size)
		for v, r := range tt.vars {
			bn.SetParents(v, r.Data().(varset.Varset), r.Score())
		}
		for v := 0; v < tt.size; v++ {
			ps := bn.Parents(v)
			r, ok := tt.vars[v]
			if ok {
				if !r.Data().(varset.Varset).Equal(ps) {
					t.Errorf("wrong equality of (%v): (%v)!=(%v)", v, r.Data().(varset.Varset), ps)
				}
			} else {
				if nil != ps {
					t.Errorf("wrong equality of (%v): (%v)!=(%v)", v, nil, ps)
				}
			}
		}
	}
}

func TestBetter(t *testing.T) {
	cases := []struct {
		size         int
		vars1, vars2 map[int]*scr.Record
		better       bool
	}{}
	for _, tt := range cases {
		bn1, bn2 := NewBNStructure(tt.size), NewBNStructure(tt.size)
		for v, r := range tt.vars1 {
			bn1.SetParents(v, r.Data().(varset.Varset), r.Score())
		}
		for v, r := range tt.vars2 {
			bn2.SetParents(v, r.Data().(varset.Varset), r.Score())
		}
		if bn1.Better(bn2) != tt.better {
			t.Errorf("wrong compare between\n(%v)\n(%v)\n(%v)!=(%v)", bn1, bn2, bn1.Better(bn2), tt.better)
		}
	}
}

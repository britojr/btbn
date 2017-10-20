package scr

import (
	"os"
	"testing"

	"github.com/britojr/utl/ioutl"
)

func TestComputeFromDataset(t *testing.T) {
	cases := []struct {
		content string
		nvar    int
		mi      [][]float64
	}{{
		`A,B,C
0,1,0
0,1,1
0,0,1
1,1,1
0,1,0`, 3, [][]float64{
			{0.50040242353, 0.05053430783, 0.11849392254},
			{0.05053430783, 0.50040242353, 0.11849392254},
			{0.11849392254, 0.11849392254, 0.673011667},
		},
	}}

	for _, tt := range cases {
		fname := ioutl.TempFile("mi_test", tt.content)
		defer os.Remove(fname)
		mi := ComputeMutInf(fname)
		if tt.nvar != mi.NVar() {
			t.Errorf("wrong var number (%v)!=(%v)", tt.nvar, mi.NVar())
		}
		for i := 0; i < mi.NVar(); i++ {
			for j := 0; j < mi.NVar(); j++ {
				got := mi.Get(i, j)
				if tt.nvar != mi.NVar() {
					t.Errorf("wrong mi[%v][%v]: (%v)!=(%v)", i, j, tt.mi[i][j], got)
				}
			}
		}
	}
}

func TestWriteRead(t *testing.T) {
	cases := []struct {
		content string
		nvar    int
		mi      [][]float64
	}{{
		`A,B,C
0,1,0
0,1,1
0,0,1
1,1,1
0,1,0`, 3, [][]float64{
			{0.50040242353, 0.05053430783, 0.11849392254},
			{0.05053430783, 0.50040242353, 0.11849392254},
			{0.11849392254, 0.11849392254, 0.673011667},
		},
	}}

	for _, tt := range cases {
		dsfile := ioutl.TempFile("mi_test", tt.content)
		mifile := ioutl.TempFile("mi_test", tt.content)
		defer os.Remove(dsfile)
		defer os.Remove(mifile)
		mi1 := ComputeMutInf(dsfile)
		mi1.Write(mifile)
		mi := ReadMutInf(mifile)
		if tt.nvar != mi.NVar() {
			t.Errorf("wrong var number (%v)!=(%v)", tt.nvar, mi.NVar())
		}
		for i := 0; i < mi.NVar(); i++ {
			for j := 0; j < mi.NVar(); j++ {
				got := mi.Get(i, j)
				if tt.nvar != mi.NVar() {
					t.Errorf("wrong mi[%v][%v]: (%v)!=(%v)", i, j, tt.mi[i][j], got)
				}
			}
		}
	}
}

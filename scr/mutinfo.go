package scr

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"math"
	"strings"

	"github.com/britojr/utl/conv"
	"github.com/britojr/utl/errchk"
	"github.com/britojr/utl/ioutl"
)

// MutInfo handles mutual information
type MutInfo struct {
	mat [][]float64
}

// Get returns the mutual information of a given pair of variables
func (m *MutInfo) Get(i, j int) float64 {
	if i > j {
		return m.mat[i][j]
	}
	return m.mat[j][i]
}

// NVar returns the number of variables
func (m *MutInfo) NVar() int {
	return len(m.mat)
}

// Write writes mutual information on file
func (m *MutInfo) Write(fname string) {
	f := ioutl.CreateFile(fname)
	defer f.Close()
	for i := range m.mat {
		for j := 0; j < i; j++ {
			fmt.Fprintf(f, "%v ", m.mat[i][j])
		}
		fmt.Fprintf(f, "%v\n", m.mat[i][i])
	}
}

// ReadMutInfo reads a mutual information file
func ReadMutInfo(fname string) *MutInfo {
	f := ioutl.OpenFile(fname)
	defer f.Close()
	scanner := bufio.NewScanner(f)

	mat := [][]float64(nil)
	for scanner.Scan() {
		line := conv.Satof(strings.Fields(scanner.Text()))
		mat = append(mat, line)
	}
	m := new(MutInfo)
	m.mat = mat
	return m
}

// ComputeFromDataset reads a dataset and computes mutual information
func ComputeFromDataset(fname string) *MutInfo {
	// TODO: split in 3 steps: counting, entropy and information
	f := ioutl.OpenFile(fname)
	defer f.Close()
	r := csv.NewReader(bufio.NewReader(f))

	var (
		N   int
		mx  [][]int     // individual count
		mxy [][][][]int // pair count
		mat [][]float64 // pair mutual information
	)
	r.Read() // ignore header line
	record, err := r.Read()
	for ; err != io.EOF; record, err = r.Read() {
		errchk.Check(err, "")
		N++
		line := conv.Satoi(record)
		if N == 1 {
			mat = make([][]float64, len(line))
			mx = make([][]int, len(line))
			mxy = make([][][][]int, len(line))
			for i := range line {
				mxy[i] = make([][][]int, i)
				mat[i] = make([]float64, i+1)
			}
		}
		for i := range line {
			for len(mx[i]) <= line[i] {
				mx[i] = append(mx[i], 0)
			}
			mx[i][line[i]]++
			for j := 0; j < i; j++ {
				for len(mxy[i][j]) <= line[i] {
					mxy[i][j] = append(mxy[i][j], []int{})
				}
				for len(mxy[i][j][line[i]]) <= line[j] {
					mxy[i][j][line[i]] = append(mxy[i][j][line[i]], 0)
				}
				mxy[i][j][line[i]][line[j]]++
			}
		}
	}
	logN := math.Log(float64(N))
	hx := make([]float64, len(mx))
	for i := range mx {
		for k := range mx[i] {
			if mx[i][k] > 0 {
				hx[i] += float64(mx[i][k]) * math.Log(float64(mx[i][k]))
			}
		}
	}
	for i := range mxy {
		for j := 0; j < i; j++ {
			hij := float64(0)
			for u := range mxy[i][j] {
				for v := range mxy[i][j][u] {
					if mxy[i][j][u][v] > 0 {
						hij += float64(mxy[i][j][u][v]) * math.Log(float64(mxy[i][j][u][v]))
					}
				}
			}
			mat[i][j] = ((hij - hx[i] - hx[j]) / float64(N)) + logN
		}
		mat[i][i] = (-hx[i] / float64(N)) + logN
	}

	m := new(MutInfo)
	m.mat = mat
	return m
}

package scr

import (
	"bufio"
	"strings"

	"github.com/britojr/utl/conv"
	"github.com/britojr/utl/ioutl"
)

// MutInfo handles mutual information
type MutInfo struct {
	mat [][]float64
}

// Get returns the mutual information of a given pair of variables
func (m MutInfo) Get(i, j int) float64 {
	if i > j {
		return m.mat[i][j]
	}
	return m.mat[j][i]
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

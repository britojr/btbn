package scr

import (
	"bufio"
	"fmt"
	"strings"

	"gonum.org/v1/gonum/stat"

	"github.com/britojr/btbn/dfext"
	"github.com/britojr/utl/conv"
	"github.com/britojr/utl/ioutl"
	"github.com/kniren/gota/dataframe"
)

// MutInfo handles mutual information
type MutInfo struct {
	mat [][]float64 // lower triangular matrix with pairwise mutual information
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

// ReadMutInf reads a mutual information file
func ReadMutInf(fname string) *MutInfo {
	f := ioutl.OpenFile(fname)
	defer f.Close()
	scanner := bufio.NewScanner(f)

	mat := [][]float64(nil)
	for scanner.Scan() {
		line := conv.Satof(strings.Fields(scanner.Text()))
		mat = append(mat, line)
	}
	return &MutInfo{mat}
}

// ComputeMutInf computes mutual information from a csv dataset
func ComputeMutInf(fname string) *MutInfo {
	f := ioutl.OpenFile(fname)
	defer f.Close()
	df := dataframe.ReadCSV(bufio.NewReader(f))
	return ComputeMutInfDF(df)
}

// ComputeMutInfDF computes mutual information from a dataframe
func ComputeMutInfDF(df dataframe.DataFrame) *MutInfo {
	mat := make([][]float64, df.Ncol())
	for i := range mat {
		mat[i] = make([]float64, i+1)
	}
	// compute empiric individual entropy for each variable and store it in the diagonal
	for i := 0; i < df.Ncol(); i++ {
		mat[i][i] = stat.Entropy(dfext.Counts(df.Select([]int{i}), true))
	}
	// compute pairwise mutual information for each pair and store in matrix lower triangle
	for i := 0; i < df.Ncol(); i++ {
		for j := 0; j < i; j++ {
			mat[i][j] = stat.Entropy(dfext.Counts(df.Select([]int{i, j}), true))
		}
	}
	return &MutInfo{mat}
}

// Package scr provides methods to handle score information
package scr

import (
	"bufio"
	"strings"

	"github.com/britojr/btbn/varset"
	"github.com/britojr/utl/conv"
	"github.com/britojr/utl/ioutl"
)

// Cache stores pre-computed score information
type Cache struct {
	nvar     int
	caches   []map[string]float64
	varName  []string
	varIndex map[string]int
}

// Read reads a score file into a score cache
func Read(fname string) *Cache {
	f := ioutl.OpenFile(fname)
	defer f.Close()
	scanner := bufio.NewScanner(f)

	// read all variables
	varIndex := make(map[string]int)
	varName := []string(nil)
	for scanner.Scan() {
		words := strings.Fields(scanner.Text())
		if len(words) != 0 && words[0] == "VAR" {
			varName = append(varName, words[1])
			varIndex[words[1]] = len(varName) - 1
		}
	}

	c := &Cache{}
	c.nvar = len(varName)
	c.varIndex = varIndex
	c.varName = varName
	c.caches = make([]map[string]float64, c.Nvar())
	for i := range c.caches {
		c.caches[i] = make(map[string]float64)
	}

	// rewind file to read parents score
	f.Seek(0, 0)
	scanner = bufio.NewScanner(f)
	currVar := 0
	for scanner.Scan() {
		words := strings.Fields(scanner.Text())
		if len(words) == 0 || words[0] == "META" {
			continue
		}
		if words[0] == "VAR" {
			currVar = c.varIndex[words[1]]
			continue
		}
		scoreVal := conv.Atof(words[0])

		parents := varset.New(c.Nvar())
		for i := 1; i < len(words); i++ {
			parents.Set(c.varIndex[words[i]])
		}
		c.putScore(currVar, parents, scoreVal)
	}
	return c
}

func (c *Cache) putScore(v int, parents varset.Varset, scoreVal float64) {
	c.caches[v][parents.DumpHashString()] = scoreVal
}

// Nvar returns the number of variables
func (c *Cache) Nvar() int {
	return c.nvar
}

// Scores returns the score map for a variable v
func (c *Cache) Scores(v int) map[string]float64 {
	return c.caches[v]
}

// VarName returns the name of a given variable id
func (c *Cache) VarName(v int) string {
	return c.varName[v]
}

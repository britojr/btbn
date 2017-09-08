package score

import (
	"bufio"
	"log"
	"strconv"
	"strings"

	"github.com/britojr/btbn/varset"
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
		scoreVal, err := strconv.ParseFloat(words[0], 64)
		if err != nil {
			log.Printf("Error trying to convert %v to float\n", words[0])
			panic(err)
		}

		parents := varset.New(c.Nvar())
		for i := 1; i < len(words); i++ {
			parents.Set(c.varIndex[words[i]])
		}
		c.putScore(currVar, parents, scoreVal)
	}
	return c
}

func (c *Cache) putScore(v int, parents varset.Varset, scoreVal float64) {
	c.caches[v][parents.HashString()] = scoreVal
}

// Nvar returns the number of variables
func (c *Cache) Nvar() int {
	return c.nvar
}

// Scores returns the score map for a variable v
func (c *Cache) Scores(v int) map[string]float64 {
	return c.caches[v]
}

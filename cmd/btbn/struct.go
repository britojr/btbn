package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"time"

	"github.com/britojr/btbn/optimizer"
	"github.com/britojr/btbn/scr"
	"github.com/britojr/btbn/varset"
	"github.com/britojr/utl/ioutl"
)

func runStructComm() {
	// Required Flags
	if scoreFile == "" {
		fmt.Printf("\n error: missing score file\n\n")
		structComm.PrintDefaults()
		os.Exit(1)
	}
	if !verbose {
		log.SetOutput(ioutil.Discard)
	}

	structureLearning()
}

func structureLearning() {
	log.Printf(" ========== BEGIN STRUCTURE OPTIMIZATION ========== \n")
	log.Printf("Learning algorithm: '%v'\n", optimizerAlg)
	log.Printf("Max. iterations: %v\n", numSolutions)
	log.Printf("Max. time available (sec): %v\n", timeAvailable)
	log.Printf("Pre-computed scores file: '%v'\n", scoreFile)
	log.Printf("Parameters file: '%v'\n", parmFile)
	log.Printf("Save solution in: '%v'\n", bnetFile)
	log.Printf(" -------------------------------------------------- \n")

	log.Println("Reading parameters file")
	parms := ioutl.ReadYaml(parmFile)
	maxPa := defMaxPa(parms)
	log.Println("Reading score cache")
	scoreCache := scr.Read(scoreFile)
	log.Println("Creating score rankers")
	scoreRankers := scr.CreateRankers(scoreCache, maxPa)
	// TODO: dataset will also be nedded when dealing with hidden variables

	log.Println("Creating bounded-treewidth structure learning algorithm")
	algorithm := optimizer.Create(optimizerAlg, scoreRankers, parms)
	algorithm.PrintParameters()

	log.Println("Searching bounded-treewidth structure")
	start := time.Now()
	solution := optimizer.Search(algorithm, numSolutions, timeAvailable)
	elapsed := time.Since(start)

	totScore := solution.Score()
	empScore := emptySetScore(scoreRankers)
	log.Printf(" ========== SOLUTION ============================ \n")
	if solution == nil {
		log.Printf("Couldn't find any solution in the given time!\n")
	} else {
		log.Printf("Time: %v\n", elapsed)
		log.Printf("Best Score: %.6f\n", totScore)
		log.Printf("Normalized: %.6f\n", (totScore-empScore)/math.Abs(empScore))
	}
	log.Printf(" -------------------------------------------------- \n")

	if len(bnetFile) > 0 {
		writeSolution(bnetFile, solution)
	}
}

func writeSolution(fname string, bn *optimizer.BNStructure) {
	// bn.Size() bn.Parents(i).DumpString()
	log.Printf("Printing solution: '%v'\n", fname)
	f := ioutl.CreateFile(fname)
	defer f.Close()
	fmt.Fprintf(f, "META variables = %v\n", bn.Size())
	fmt.Fprintf(f, "META score = %v\n", bn.Score())
	for i := 0; i < bn.Size(); i++ {
		fmt.Fprintf(f, "%v:", i)
		for _, v := range bn.Parents(i).DumpAsInts() {
			fmt.Fprintf(f, " %v", v)
		}
		fmt.Fprintf(f, "\n")
	}
}

// emptySetScore calculates the total score for when the parents sets are empty
func emptySetScore(rankers []scr.Ranker) (es float64) {
	parents := varset.New(len(rankers))
	for _, ranker := range rankers {
		es += ranker.ScoreOf(parents)
	}
	return
}

func defMaxPa(parms map[string]string) int {
	return 0
}

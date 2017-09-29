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
	"github.com/britojr/utl/conv"
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
	maxPa := getMaxPa(parms)
	log.Printf("%v: '%v'\n", optimizer.ParmMaxParents, maxPa)
	log.Println("Reading pre-computed scores file")
	scoreCache := scr.Read(scoreFile)
	scoreRanker := scr.CreateRanker(scoreCache, maxPa)
	// TODO: dataset will also be nedded when dealing with hidden variables

	log.Println("Creating bounded-treewidth structure learning algorithm")
	algorithm := optimizer.Create(optimizerAlg, scoreRanker, parms)
	algorithm.PrintParameters()

	log.Println("Searching bounded-treewidth structure")
	start := time.Now()
	solution := optimizer.Search(algorithm, numSolutions, timeAvailable)
	elapsed := time.Since(start)

	log.Printf(" ========== SOLUTION ============================ \n")
	if solution == nil {
		log.Printf("Couldn't find any solution in the given time!\n")
		os.Exit(0)
	}
	totScore := solution.Score()
	empScore := emptySetScore(scoreRanker)
	log.Printf("Time: %v\n", elapsed)
	log.Printf("Best Score: %.6f\n", totScore)
	log.Printf("Normalized: %.6f\n", (totScore-empScore)/math.Abs(empScore))
	log.Printf(" -------------------------------------------------- \n")

	if len(bnetFile) > 0 {
		writeSolution(bnetFile, solution, algorithm, scoreCache)
	}
}

func writeSolution(fname string, bn *optimizer.BNStructure, alg optimizer.Optimizer, sc *scr.Cache) {
	log.Printf("Printing solution: '%v'\n", fname)
	f := ioutl.CreateFile(fname)
	defer f.Close()
	fmt.Fprintf(f, "META variables = %v\n", bn.Size())
	fmt.Fprintf(f, "META treewidth = %v\n", alg.Treewidth())
	fmt.Fprintf(f, "META score = %v\n", bn.Score())
	fmt.Fprintln(f)
	for i := 0; i < bn.Size(); i++ {
		fmt.Fprintf(f, "%v:", sc.VarName(i))
		for _, v := range bn.Parents(i).DumpAsInts() {
			fmt.Fprintf(f, " %v", sc.VarName(v))
		}
		fmt.Fprintf(f, "\n")
	}
}

// emptySetScore calculates the total score for when the parents sets are empty
func emptySetScore(ranker scr.Ranker) (es float64) {
	parents := varset.New(ranker.Size())
	for v := 0; v < ranker.Size(); v++ {
		es += ranker.ScoreOf(v, parents)
	}
	return
}

func getMaxPa(parms map[string]string) int {
	var tw, mp int
	if stw, ok := parms[optimizer.ParmTreewidth]; ok {
		tw = conv.Atoi(stw)
	}
	if smp, ok := parms[optimizer.ParmMaxParents]; ok {
		mp = conv.Atoi(smp)
	}
	switch {
	case tw <= 0:
		return mp
	case mp <= 0 || mp > tw:
		return tw
	}
	return mp
}

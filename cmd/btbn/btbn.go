package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/britojr/btbn/optimizer"
	"github.com/britojr/btbn/score"
)

// Define subcommand names
const (
	structConst = "struct"
)

// Define Flag variables
var (
	// common
	verbose bool // verbose mode

	// struct command
	scoreFile     string // scores input file
	bnetFile      string // network output file
	parmFile      string // parameters file for search algorithms
	optimizerAlg  string // structure optimizer algorithm
	maxPa         int    // max parents
	timeAvailable int    // time available to search solution
	numSolutions  int    // number of iterations

	// Define subcommands
	structComm *flag.FlagSet
)

func main() {
	initSubcommands()
	// Verify that a subcommand has been provided
	// os.Arg[0] : main command, os.Arg[1] : subcommand
	if len(os.Args) < 2 {
		printDefaults()
		os.Exit(1)
	}
	switch os.Args[1] {
	case structConst:
		structComm.Parse(os.Args[2:])
		runStructComm()
	default:
		printDefaults()
		os.Exit(1)
	}
}

func initSubcommands() {
	// Subcommands
	structComm = flag.NewFlagSet(structConst, flag.ExitOnError)

	// struct subcommand flags
	structComm.BoolVar(&verbose, "v", false, "prints detailed steps")

	structComm.StringVar(&scoreFile, "s", "", "precomputed scores file")
	structComm.StringVar(&parmFile, "p", "", "parameters file")
	structComm.StringVar(&bnetFile, "b", "", "network output file")

	structComm.StringVar(&optimizerAlg, "a", "sample", "structure optimizer algorithm {sample|iterative}")
	structComm.IntVar(&timeAvailable, "t", 60, "available time to search solution (0->unbounded)")
	structComm.IntVar(&numSolutions, "i", 1, "max number of iterations (0->unbounded)")
	structComm.IntVar(&maxPa, "mp", 0, "max number of parents (0->unbounded)")
}

func printDefaults() {
	fmt.Printf("Usage:\n\n")
	fmt.Printf("\tbtbn <command> [arguments]\n\n")
	fmt.Printf("The commands are:\n\n")
	fmt.Printf("\t%v\n",
		structConst,
	)
	fmt.Println()
}

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
	log.Printf("========== STEP: STRUCTURE OPTIMIZATION ========== \n")
	log.Printf("Learning algorithm: '%v'\n", optimizerAlg)
	log.Printf("Max. iterations: %v\n", numSolutions)
	log.Printf("Max. time available (sec): %v\n", timeAvailable)
	log.Printf("Pre-computed scores file: '%v'\n", scoreFile)
	log.Printf("Parameters file: '%v'\n", parmFile)
	log.Printf("Save solution in: '%v'\n", bnetFile)

	log.Println("Reading score cache")
	scoreCache := score.Read(scoreFile)

	log.Println("Creating score rankers")
	// scoreRankers := score.CreateRankers(scoreRankType, scoreCache, maxPa)
	scoreRankers := score.CreateRankers(scoreCache, maxPa)
	// TODO: may need to load a pre-computed mutual information file
	// TODO: dataset will also be nedded when dealing with hidden variables

	log.Println("Creating bounded-treewidth structure learning algorithm")
	algorithm := optimizer.Create(optimizerAlg, scoreRankers, parmFile)
	// algorithm.PrintParameters()

	log.Println("Searching bounded-treewidth structure")
	solution := optimizer.Search(algorithm, numSolutions, timeAvailable)
	writeSolution(bnetFile, solution)
}

func writeSolution(fname string, bnet *optimizer.BNStructure) {
	// datastructures::BNStructure bnet
	// log.Printf("Time: %v, Total score: %v, Normalized: %v\n", elapsed,
	// -solution.getScore(), scoreFunction.Normalize(bestScore))
	// writeOutput(resultFile,
	// "tree-with,norm-score,num-var,iterations,elapsed-time\n",
	// k, scoreFunction.Normalize(bestScore), n, iterations, elapsed)
	log.Printf("Printing solution: '%v'\n", fname)
}

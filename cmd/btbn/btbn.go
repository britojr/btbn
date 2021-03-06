package main

import (
	"flag"
	"fmt"
	"os"
)

// Define subcommand names
const (
	structConst = "struct"
	structDescr = "run bounded tree-width structure learning algorithm"
	mutinfConst = "mutinf"
	mutinfDescr = "computes pairwise mutual information"
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
	timeAvailable int    // time available to search solution
	numSolutions  int    // number of iterations

	// mutinf command
	dataFile   string // dataset csv file
	mutinfFile string // mutual information file

	// Define subcommands
	structComm *flag.FlagSet
	mutinfComm *flag.FlagSet
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
	case mutinfConst:
		mutinfComm.Parse(os.Args[2:])
		runMutinfComm()
	default:
		printDefaults()
		os.Exit(1)
	}
}

func initSubcommands() {
	// Subcommands
	structComm = flag.NewFlagSet(structConst, flag.ExitOnError)
	mutinfComm = flag.NewFlagSet(mutinfConst, flag.ExitOnError)

	// struct subcommand flags
	structComm.BoolVar(&verbose, "v", true, "prints detailed steps")
	structComm.StringVar(&scoreFile, "s", "", "precomputed scores file")
	structComm.StringVar(&parmFile, "p", "", "parameters file")
	structComm.StringVar(&bnetFile, "b", "", "network output file")
	structComm.StringVar(&optimizerAlg, "a", "sample", "structure optimizer algorithm {sample|selected|guided|iterative}")
	structComm.IntVar(&timeAvailable, "t", 60, "available time to search solution (0->unbounded)")
	structComm.IntVar(&numSolutions, "i", 1, "max number of iterations (0->unbounded)")

	// mutinf subcommand Flags
	mutinfComm.BoolVar(&verbose, "v", true, "prints detailed steps")
	mutinfComm.StringVar(&dataFile, "d", "", "dataset file in csv format")
	mutinfComm.StringVar(&mutinfFile, "o", "", "file to save mutual information")
}

func printDefaults() {
	fmt.Printf("btbn is a tool for learning bounded tree-width Bayesian networks\n")
	fmt.Printf("Usage:\n\n")
	fmt.Printf("\tbtbn <command> [options]\n\n")
	fmt.Printf("Commands:\n\n")
	fmt.Printf("\t%v\t\t%v\n", structConst, structDescr)
	fmt.Printf("\t%v\t\t%v\n", mutinfConst, mutinfDescr)
	fmt.Println()
	fmt.Printf("For usage details of each command, run:\n\n")
	fmt.Printf("\tbtbn <command> --help\n")
	fmt.Println()
}

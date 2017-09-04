package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
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
	k       int    // treewidth
	p       int    // max parents
	nk      int    // number of k-trees to sample
	scorefi string // scores file
	ctreefo string // cliquetree output file

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
	structComm.StringVar(&scorefi, "si", "", "scores input file")
	structComm.StringVar(&ctreefo, "co", "", "cliquetree output file")
	structComm.IntVar(&k, "k", 3, "treewidth of the structure")
	structComm.IntVar(&p, "p", 3, "max number of parents")
	structComm.IntVar(&nk, "nk", 1, "number of ktrees samples")
	structComm.BoolVar(&verbose, "v", false, "prints detailed steps")
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
	if scorefi == "" {
		fmt.Printf("\n error: missing score file\n\n")
		structComm.PrintDefaults()
		os.Exit(1)
	}
	if !verbose {
		log.SetFlags(0)
		log.SetOutput(ioutil.Discard)
	}
	// print arguments
	// ...
	// run command
}

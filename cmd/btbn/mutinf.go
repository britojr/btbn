package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/britojr/btbn/scr"
)

func runMutinfComm() {
	// Required Flags
	if dataFile == "" {
		fmt.Printf("\n error: missing dataset file\n\n")
		mutinfComm.PrintDefaults()
		os.Exit(1)
	}
	if mutinfFile == "" {
		fmt.Printf("\n error: missing output file\n\n")
		mutinfComm.PrintDefaults()
		os.Exit(1)
	}
	if !verbose {
		log.SetOutput(ioutil.Discard)
	}

	mutinfComputing()
}

func mutinfComputing() {
	log.Printf(" ========== COMPUTING MUTUAL INFORMATION ========== \n")
	log.Printf("Dataset file: '%v'\n", dataFile)
	log.Printf("Save values in: '%v'\n", mutinfFile)
	log.Printf(" -------------------------------------------------- \n")

	log.Println("Computing...")
	start := time.Now()
	mi := scr.ComputeMutInf(dataFile)
	log.Println("Saving...")
	mi.Write(mutinfFile)
	elapsed := time.Since(start)
	log.Printf(" ========== DONE ============================ \n")
	log.Printf("Time: %v\n", elapsed)
	log.Printf(" -------------------------------------------------- \n")
}

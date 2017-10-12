# Bounded Tree-width Bayesian Networks Learner

btbn is a tool for learning bounded tree-width Bayesian networks

[![Build Status](https://travis-ci.org/britojr/btbn.svg?branch=master)](https://travis-ci.org/britojr/btbn)
[![Coverage Status](https://coveralls.io/repos/github/britojr/btbn/badge.svg?branch=master)](https://coveralls.io/github/britojr/btbn?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/britojr/btbn)](https://goreportcard.com/report/github.com/britojr/btbn)
[![GoDoc](https://godoc.org/github.com/britojr/btbn?status.svg)](http://godoc.org/github.com/britojr/btbn)

___

## Installation and usage

### Get, install and test:

		go get -u github.com/britojr/btbn...
		go install github.com/britojr/btbn...
		go test github.com/britojr/btbn... -cover

### Usage:

		btbn --help
		Usage: btbn <command> [options]

		Commands:

			struct		run bounded tree-width structure learning algorithm
			mutinf		compute pairwise mutual information

		Usage of struct command:

			-a			structure learning algorithm {sample|selected|guided|iterative}
			-b			network output file (ex: example.bnet)
			-i			max number of iterations (0 means unbounded) (default 1)
			-p			parameters file (ex: example-parms.yaml)
			-s			precomputed parent set scores file (ex: example.pss)
			-t			available time to search for solution (0 means unbounded)
			-v=true		prints detailed steps (default true)

		Usage of mutinf command:

			-d			dataset file (ex: example.csv)
			-o			mutual information output file (ex: example.mi)
			-v=true		prints detailed steps (default true)

### Examples:

Learn bounded tree-width bn using iterative greedy search for 2 secs.
Save resulting networnk in 'example.bnet':

	cd examples/
	btbn struct -i 0 -t 2 -s example.pss -p example-parms.yaml -a iterative -b example.bnet

Compute mutual information:

	btbn mutinf -d example.csv -o example.mi

#### Parameters file:

YAML file containing specific parameters for the structure learning algorithms

##### Parameters file fields:

	treewidth:		structure's tree-width upper bound 
	num_trees:		number of ktrees to sample and sort (used on selected strategy)
	mutual_info:		precomputed mutual information file (used on selected strategy)
	search_variation:	iterative search type (greedy or astar)

##### Parameters file example:

	./examples/example-parms.yaml
	treewidth: 4
	num_trees: 100
	mutual_info: example.mi
	search_variation: greedy

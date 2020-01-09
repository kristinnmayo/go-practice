package main

import (
	"os"

	"github.com/go-practice/cmd"
	"github.com/go-practice/target"
)

func main() {
	path := os.Args[1]              // todo -- command line options
	targets := target.Walk(path)    // get list of target files
	hits := cmd.Tryscanner(targets) // search for vulnerable strings within targets

	for _, h := range hits {
		h.Display()
	}

	for _, t := range targets {
		t.Display()
	}

	cmd.Tryast()
	cmd.Tryanalysis()
}

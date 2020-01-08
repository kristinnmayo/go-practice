package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/go-practice/cmd"
	"github.com/go-practice/target"
)

// accept a path and return all filenames in tree
func walk(dir string) (targets []target.Target) {
	// func Walk(root string, walkFn WalkFunc) error
	// type WalkFunc func(path string, info os.FileInfo, err error) error
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatal(path, err)
		}
		// only check .go files
		if filepath.Ext(path) != ".go" {
			return nil
		}
		targets = append(targets, target.New(path, info.Mode().Perm()))
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	return targets
}

func main() {
	path := os.Args[1]              // todo -- command line options
	targets := walk(path)           // get list of target files
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

package main

import (
	"bufio"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-practice/cmd"
	"github.com/go-practice/hit"
	"github.com/go-practice/target"
)

// search a list of files for vulnerable strings and return list of hits
func seek(targets []target.Target) []hit.Hit {
	// list vulnerable strings to search for
	hitlist := []string{"Sprintf", "todo", "Mkdir", "MkdirAll"}

	// slice to store hits
	var hits []hit.Hit

	for _, target := range targets {
		// open in read-only mode -> returns pointer of type os.File
		f, err := os.Open(target.Path)
		if err != nil {
			log.Fatalf("failed opening file: %s", err)
		}
		defer f.Close()

		scanner := bufio.NewScanner(f)
		scanner.Split(bufio.ScanLines)

		// Scan() forwards to the next line
		for line := 1; scanner.Scan(); line++ {
			code := scanner.Text()
			for _, vuln := range hitlist {
				if strings.Contains(code, vuln) {
					h := hit.New(target.Path, code, vuln, line)
					target.Vulns[vuln]++
					hits = append(hits, h)
				}
			}
		}
	}
	return hits
}

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
	path := os.Args[1]    // todo -- command line options
	targets := walk(path) // get list of target files
	hits := seek(targets) // search for vulnerable strings within targets

	for _, h := range hits {
		h.Display()
	}

	for _, t := range targets {
		t.Display()
	}

	cmd.Tryscanner()
	cmd.Tryast()
	cmd.Tryanalysis()
}

package main

import (
	"bufio"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// Hit is a possible vulnerability
type Hit struct {
	file string
	line int
	code string
	vuln string
	err  error
}

// log info on a hit
func (h Hit) display() {
	log.Printf("%s:%d %s (%s)", h.file, h.line, h.code, h.vuln)
}

// search a list of files for vulnerable strings and return list of hits
func seek(filenames []string) (hits []Hit) {
	// list vulnerable strings to search for
	hitlist := []string{"Sprintf", "todo", "Mkdir", "MkdirAll"}

	for _, file := range filenames {
		// open in read-only mode -> returns pointer of type os.File
		f, err := os.Open(file)
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
					h := Hit{file, line, code, vuln, nil}
					hits = append(hits, h)
				}
			}
		}
	}
	return hits
}

// accept a path and return all filenames in tree
func walk(dir string) (files []string) {
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		// only check .go files
		if filepath.Ext(path) != ".go" {
			return nil
		}
		files = append(files, path)
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	return files
}

func main() {
	// command line options
	path := ""
	args := os.Args
	if len(args) > 1 {
		path = args[1]
	} else {
		log.Fatal("\n\tsyntax: go-tools <filename.go>\n")
	}

	// get list of target files
	targets := walk(path)

	// search for vulnerable strings within targets
	hits := seek(targets)

	// display details of each found possible vulnerability
	for _, h := range hits {
		h.display()
	}
}

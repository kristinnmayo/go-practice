package main

import (
	"bufio"
	"os"
	"log"
	"strings"
	"path/filepath"
)

type Filelist []string

type Hit struct {
	file string
	line int
	code string
	vuln string
}

// log info on a hit
func (h Hit) display() {
	log.Printf("%s:%d %s (%s)", h.file, h.line, h.code, h.vuln)
}

// search a list of files for vulnerable strings and return list of hits
func (filenames Filelist) seek() (hits []Hit) {
	// list vulnerable strings to search for
	hitlist := []string{"sprintf", "todo"}

	// for each file in tree, do stuff
	for _, file := range filenames {
		// open in read-only mode (returns pointer of type os.File)
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
					h := Hit{file, line, code, vuln}
					hits = append(hits, h)
				}
			}
		}
	}
	return hits
}

// accept a path and return all filenames in tree
func walk(path string) (files []string) {
	root := "./"

	// func Walk(root string, walkFn WalkFunc) error
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
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
		log.Fatal("\n\tsyntax: ./go-tools <filename.go>\n")
	}

	// list all paths in tree
	filenames := Filelist(walk(path))

	// search for vulnerable strings
	hits := filenames.seek()
	for _, h := range hits {
		h.display()
	}
}

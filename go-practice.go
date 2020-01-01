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
	file string // path to file
	line int    // line number
	code string // vulnerable line of code
	vuln string // type of vulnerability
	err  error
}

// File is one file in the user specified system
type File struct {
	name  string         // name of file
	perm  os.FileMode    // file permissions
	vulns map[string]int // number of hits found
}

// log info on a hit
func (h Hit) display() {
	log.Printf("%s:%d %s (%s)", h.file, h.line, h.code, h.vuln)
}

// log info on a file
func (f File) display() {
	log.Printf("%s (%04o)", f.name, f.perm)
	log.Println(f.vulns)
}

// search a list of files for vulnerable strings and return list of hits
func seek(files []File) (hits []Hit) {
	// list vulnerable strings to search for
	hitlist := []string{"Sprintf", "todo", "Mkdir", "MkdirAll"}

	for _, file := range files {
		// open in read-only mode -> returns pointer of type os.File
		f, err := os.Open(file.name)
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
					h := Hit{file.name, line, code, vuln, nil}
					file.vulns[vuln]++
					hits = append(hits, h)
				}
			}
		}
	}
	return hits
}

// accept a path and return all filenames in tree
func walk(dir string) (files []File) {
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Println(err)
		}
		// only check .go files
		if filepath.Ext(path) != ".go" {
			return nil
		}

		name := info.Name()
		perm := info.Mode().Perm()
		vulns := make(map[string]int)
		f := File{name, perm, vulns}
		files = append(files, f)
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
		log.Fatal("\n\tusage: go-tools <filename.go>\n")
	}

	// get list of target files
	targets := walk(path)

	// search for vulnerable strings within targets
	hits := seek(targets)

	// display details of each found possible vulnerability
	for _, h := range hits {
		h.display()
	}

	// display file permissions for each file
	for _, t := range targets {
		t.display()
	}
}

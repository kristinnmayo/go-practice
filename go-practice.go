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
	path  string         // file path
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

func output(hits []Hit, files []File) {
	for _, h := range hits {
		h.display()
	}
	for _, f := range files {
		f.display()
	}
}

// search a list of files for vulnerable strings and return list of hits
func seek(files []File) (hits []Hit) {
	// list vulnerable strings to search for
	hitlist := []string{"Sprintf", "todo", "Mkdir", "MkdirAll"}

	for _, file := range files {
		// open in read-only mode -> returns pointer of type os.File
		f, err := os.Open(file.path)
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

		name := info.Name()
		perm := info.Mode().Perm()
		vulns := make(map[string]int)
		f := File{name, path, perm, vulns}
		files = append(files, f)
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	return files
}

func main() {
	path := os.Args[1]    // todo -- command line options
	targets := walk(path) // get list of target files
	hits := seek(targets) // search for vulnerable strings within targets
	output(hits, targets) // print results
}

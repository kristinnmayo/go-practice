package main

import (
	"bufio"
	"os"
	"log"
	"strings"
	"path/filepath"
)

type hit struct {
	file string
	line int
	column int
	code string
	err error
}

// accept a path and return all filenames
func walk(path string) (files []string) {
	root := "./"
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
	filenames := walk(path)

	// list vulnerable strings to search for
	hitlist := []string{"sprintf", "todo", "note"}

	// for each file in tree, do stuff
	for _, file := range filenames {
		log.Println(file)

		// open in read-only mode (returns pointer of type os.File)
		f, err := os.Open(file)
		if err != nil {
			log.Fatalf("failed opening file: %s", err)
		}

		// execute after current function returns
		defer f.Close()

		scanner := bufio.NewScanner(f)
		scanner.Split(bufio.ScanLines)

		// Scan() forwards to the next line
		for l := 1; scanner.Scan(); l++ {
			line := scanner.Text()
			for _, v := range hitlist {
				if strings.Contains(line, v) {
					log.Printf("%s:%d:%s\n", file, l, line)
				}
			}
		}
	}
}

package cmd

import (
	"bufio"
	"log"
	"os"
	"strings"

	"github.com/go-practice/hit"
	"github.com/go-practice/target"
)

// Tryscanner is used to search a list of files for vulnerable strings and return list of hits
func Tryscanner(targets []target.Target) []hit.Hit {
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

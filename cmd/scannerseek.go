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
func Tryscanner(targets []*target.Target) {
	var hits []*hit.Hit
	// list vulnerable strings to search for
	hitlist := []string{"Sprintf", "todo"}
	for _, t := range targets {
		// open in read-only mode -> returns pointer of type os.File
		f, err := os.Open(t.Path)
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
					h := hit.New(t.Path, code, vuln, line)
					t.Vulns[vuln]++
					hits = append(hits, &h)
				}
			}
		}
	}
	for _, h := range hits {
		h.Display()
	}
}

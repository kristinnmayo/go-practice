package cmd

import (
	"github.com/go-practice/sprintfcheck"
	"golang.org/x/tools/go/analysis/multichecker"
)

// Analysisseek ...
func Analysisseek() {
	multichecker.Main(
		sprintfcheck.Analyzer,
	)
}

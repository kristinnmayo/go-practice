// Package todocheck defines an Analyzer that reports...
package todocheck

import (
	"fmt"
	"strings"

	"golang.org/x/tools/go/analysis"
)

// The analysis package defines the interface between a modular static analysis and an analysis driver program.
// The primary type in the API is Analyzer. An Analyzer statically describes an analysis function: its name,
// documentation, flags, relationship to other analyzers, and of course, its logic.

// Analyzer ...
var Analyzer = &analysis.Analyzer{
	Name: "todocheck",
	Doc:  "seeks leftover todos",

	// On success, the Run function may return a result
	// computed by the Analyzer; its type must match ResultType.
	// The driver makes this result available as an input to
	// another Analyzer that depends directly on this one (see
	// Requires) when it analyzes the same package.
	// To pass analysis results between packages (and thus
	// potentially between address spaces), use Facts, which are
	// serializable.
	Run: run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		for _, i := range file.Comments {
			if strings.Contains(i.Text(), "todo") {
				fmt.Printf("%s:\t%s\n", pass.Fset.Position(i.Pos()), i.Text())
			}
		}
	}

	return nil, nil
}

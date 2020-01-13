// Package sprintfcheck defines an Analyzer that reports...
package sprintfcheck

import (
	"fmt"
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

// Analyzer ...
var Analyzer = &analysis.Analyzer{
	Name: "sprintfcheck",
	Doc:  "seeks use of sprintf",

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
		ast.Inspect(file, func(n ast.Node) bool {

			var s string
			switch x := n.(type) {
			case *ast.Ident:
				s = x.Name
				if s == "Sprintf" {
					fmt.Printf("%s:\t%s\n", pass.Fset.Position(n.Pos()), s)
				}
			}
			return true
		})
	}

	return nil, nil
}

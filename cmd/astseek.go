package cmd

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"strings"

	"github.com/go-practice/types/target"
)

// Astseek todo
func Astseek(targets []*target.Target) {
	for _, t := range targets {
		fileset := token.NewFileSet() // positions are relative to fset
		root, err := parser.ParseFile(fileset, t.Path, nil, parser.ParseComments)
		if err != nil {
			panic(err)
		}

		// manual inspection
		for _, i := range root.Comments {
			if strings.Contains(i.Text(), "todo") {
				log.Println(i.Text())
				t.Vulns["todo"]++
			}
		}

		// depth first traversal
		ast.Inspect(root, func(n ast.Node) bool {
			switch x := n.(type) {
			case *ast.Ident:
				s := x.Name
				switch s {
				case "Sprintf", "Mkdir", "MkdirAll":
					fmt.Printf("%s:\t%s\n", fileset.Position(n.Pos()), s)
					// t.Vulns["Sprintf"]++
				}
			}
			return true
		})
	}
}

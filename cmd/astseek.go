package cmd

import (
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"strings"

	"github.com/go-practice/hit"
	"github.com/go-practice/target"
)

type visitor int

// Tryast ...
func Tryast() {
	// todo -- command line options
	paths := os.Args[1:]
	var targets []target.Target
	var hits []hit.Hit

	// list vulnerable strings to search for
	hitlist := []string{"Sprintf", "todo", "Mkdir", "MkdirAll"}

	// create set of source files for parser - positions are relative to fileset Mkdir
	fileset := token.NewFileSet()

	// 	type File struct {
	//     Doc        *CommentGroup   // associated documentation; or nil
	//     Package    token.Pos       // position of "package" keyword
	//     Name       *Ident          // package name
	//     Decls      []Decl          // top-level declarations; or nil
	//     Scope      *Scope          // package scope (this file only)
	//     Imports    []*ImportSpec   // imports in this file
	//     Unresolved []*Ident        // unresolved identifiers in this file
	//     Comments   []*CommentGroup // list of all comments in the source file
	// }
	var files []*ast.File

	for _, path := range paths {
		info, err := os.Stat(path)
		if err != nil {
			log.Fatal(err)
		}
		t := target.New(path, info.Mode())
		targets = append(targets, t)

		// func ParseFile(fset *token.FileSet, filename string, src interface{}, mode Mode) (f *ast.File, err error)
		// ParseFile parses the source code of a single Go source file and returns the corresponding ast.File node.
		// The source code may be provided via the filename of the source file, or via the src parameter.

		// If src != nil, ParseFile parses the source from src and the filename is only used when recording position information.
		// The type of the argument for the src parameter must be string, []byte, or io.Reader. If src == nil, ParseFile parses the file specified by filename.

		// The mode parameter controls the amount of source text parsed and other optional parser functionality.
		// Position information is recorded in the file set fset, which must not be nil.

		// If the source couldn't be read, the returned AST is nil and the error indicates the specific failure.
		// If the source was read but syntax errors were found, the result is a partial AST (with ast.Bad* nodes representing the fragments of erroneous source code).
		// // Multiple errors are returned via a scanner.ErrorList which is sorted by file position
		root, err := parser.ParseFile(fileset, path, nil, parser.ParseComments)
		if err != nil {
			log.Fatal(err)
		}

		// add this file to our list of files
		files = append(files, root)

		// Inspect traverses an AST in depth-first order: It starts by calling f(node); node must not be nil.
		// If f returns true, Inspect invokes f recursively for each of the non-nil children of node, followed by a call of f(nil).
		ast.Inspect(root, func(node ast.Node) bool {
			switch n := node.(type) {
			// identifiers
			case *ast.Ident:
				for _, vuln := range hitlist {
					if strings.Contains(n.Name, vuln) {
						hits = append(hits, hit.New(fileset.Position(n.Pos()).Filename, "", vuln, fileset.Position(n.Pos()).Line))
					}
				}
			// functiion declarations
			case *ast.FuncDecl:
				if n.Doc.Text() == "" {
					log.Printf("missing documentation for function:\t%s:%d:%s\n", fileset.Position(n.Pos()).Filename, fileset.Position(n.Pos()).Line, n.Name.Name)
				}
			}
			return true
		})

		// Walk traverses an AST in depth-first order: It starts by calling v.Visit(node); node must not be nil.
		// If the visitor w returned by v.Visit(node) is not nil, Walk is invoked recursively with visitor w for
		// each of the non-nil children of node, followed by a call of w.Visit(nil).
		// var node visitor
		// ast.Walk(node, root)

		// manual inspection
		log.Println("Imports:")
		for _, i := range root.Imports {
			log.Println(i.Path.Value)
		}

		log.Println("Comments:")
		for _, i := range root.Comments {
			if strings.Contains(i.Text(), "todo") {
				log.Println(i.Text())
			}
		}
	}
}

func (v visitor) Visit(n ast.Node) ast.Visitor {
	if n == nil {
		return nil
	}
	log.Printf("%s%T\n", strings.Repeat("\t", int(v)), n)
	return v + 1
}

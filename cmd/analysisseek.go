package cmd

import (
	"github.com/go-practice/sprintfcheck"
	"github.com/go-practice/todocheck"
	"golang.org/x/tools/go/analysis/multichecker"
)

// Package checker defines the implementation of the checker commands.
// The same code drives the multi-analysis driver, the single-analysis
// driver that is conventionally provided for convenience along with
// each analysis package, and the test driver.

// Package singlechecker defines the main function for an analysis driver
// with only a single analysis. This package makes it easy for a provider
// of an analysis package to also provide a standalone tool that runs just
// that analysis. As such:
// func main() { singlechecker.Main(findbadness.Analyzer) }

// Package multichecker defines the main function for an analysis driver
// with several analyzers. This package makes it easy for anyone to build
// an analysis tool containing just the analyzers they need.

// checker utilizes token, parser, ast packages abstracting away most of
// the actions taken in the astseek version of this code
// see source for details:
// https://github.com/golang/tools/blob/master/go/analysis/internal/checker/checker.go#L346

// Analysisseek ...
func Analysisseek() {
	// multichecker
	multichecker.Main(
		sprintfcheck.Analyzer,
		todocheck.Analyzer,
	)
}

// An analysis driver is a program such as vet that runs a set of analyses and
// prints the diagnostics that they report. The driver program must import the
// list of Analyzers it needs. Typically each Analyzer resides in a separate
// package. To add a new Analyzer to an existing driver, add another item to the list:

// var analyses = []*analysis.Analyzer{
// 	unusedresult.Analyzer,
// 	nilness.Analyzer,
// 	printf.Analyzer,
// }

// A driver may use the name, flags, and documentation to provide on-line help
// that describes the analyses it performs. The doc comment contains a brief
// one-line summary, optionally followed by paragraphs of explanation.

// type Analyzer struct {
// 	Name             string
// 	Doc              string
// 	Flags            flag.FlagSet
// 	Run              func(*Pass) (interface{}, error)
// 	RunDespiteErrors bool
// 	ResultType       reflect.Type
// 	Requires         []*Analyzer
// 	FactTypes        []Fact
// }

// The Flags field declares a set of named (global) flag variables that control
// analysis behavior. Unlike vet, analysis flags are not declared directly in
// the command line FlagSet; it is up to the driver to set the flag variables. A
// driver for a single analysis, a, might expose its flag f directly on the command
// line as -f, whereas a driver for multiple analyses might prefix the flag name by
// the analysis name (-a.f) to avoid ambiguity. An IDE might expose the flags through
// a graphical interface, and a batch pipeline might configure them from a config
// file. See the "findcall" analyzer for an example of flags in action.

// The RunDespiteErrors flag indicates whether the analysis is equipped to handle
// ill-typed code. If not, the driver will skip the analysis if there were parse or
// type errors. The optional ResultType field specifies the type of the result value
// computed by this analysis and made available to other analyses. The Requires field
// specifies a list of analyses upon which this one depends and whose results it may
// access, and it constrains the order in which a driver may run analyses. The
// FactTypes field is discussed in the section on Modularity. The analysis package
// provides a Validate function to perform basic sanity checks on an Analyzer, such as
// that its Requires graph is acyclic, its fact and result types are unique, and so on.

// Finally, the Run field contains a function to be called by the driver to execute the
// analysis on a single package. The driver passes it an instance of the Pass type.

// A Pass describes a single unit of work: the application of a particular Analyzer to a
// particular package of Go code. The Pass provides information to the Analyzer's Run
// function about the package being analyzed, and provides operations to the Run function
// for reporting diagnostics and other information back to the driver.

// type Pass struct {
// 	Fset       *token.FileSet
// 	Files      []*ast.File
// 	OtherFiles []string
// 	Pkg        *types.Package
// 	TypesInfo  *types.Info
// 	ResultOf   map[*Analyzer]interface{}
// 	Report     func(Diagnostic)
// 	...
// }

// The Fset, Files, Pkg, and TypesInfo fields provide the syntax trees, type information,
// and source positions for a single package of Go code.

// The OtherFiles field provides the names, but not the contents, of non-Go files such as
// assembly that are part of this package. See the "asmdecl" or "buildtags" analyzers for
// examples of loading non-Go files and reporting diagnostics against them.

// The ResultOf field provides the results computed by the analyzers required by this one,
// as expressed in its Analyzer.Requires field. The driver runs the required analyzers first
// and makes their results available in this map. Each Analyzer must return a value of the
// type described in its Analyzer.ResultType field. For example, the "ctrlflow" analyzer
// returns a *ctrlflow.CFGs, which provides a control-flow graph for each function in the
// package (see golang.org/x/tools/go/cfg); the "inspect" analyzer returns a value that
// enables other Analyzers to traverse the syntax trees of the package more efficiently;
// and the "buildssa" analyzer constructs an SSA-form intermediate representation. Each of
// these Analyzers extends the capabilities of later Analyzers without adding a dependency
// to the core API, so an analysis tool pays only for the extensions it needs.

// Analyzers are provided in the form of packages that a driver program is expected to import.
// The vet command imports a set of several analyzers, but users may wish to define their own
// analysis commands that perform additional checks. To simplify the task of creating an analysis
// command, either for a single analyzer or for a whole suite, we provide the singlechecker and
// multichecker subpackages.

// type Pass: A Pass provides information to the Run function that applies a specific analyzer
// to a single Go package. It forms the interface between the analysis logic and the driver
// program, and has both input and an output components. As in a compiler, one pass may depend
// on the result computed by another.
// type Pass struct {
// 	Analyzer *Analyzer // the identity of the current analyzer

// 	// syntax and type information
// 	Fset       *token.FileSet // file position information
// 	Files      []*ast.File    // the abstract syntax tree of each file
// 	OtherFiles []string       // names of non-Go files of this package
// 	Pkg        *types.Package // type information about the package
// 	TypesInfo  *types.Info    // type information about the syntax trees
// 	TypesSizes types.Sizes    // function for computing sizes of types

// 	// Report reports a Diagnostic, a finding about a specific location
// 	// in the analyzed source code such as a potential mistake.
// 	// It may be called by the Run function.
// 	Report func(Diagnostic)

// 	// ResultOf provides the inputs to this analysis pass, which are
// 	// the corresponding results of its prerequisite analyzers.
// 	// The map keys are the elements of Analysis.Required,
// 	// and the type of each corresponding value is the required
// 	// analysis's ResultType.
// 	ResultOf map[*Analyzer]interface{}

// 	// ImportObjectFact retrieves a fact associated with obj.
// 	// Given a value ptr of type *T, where *T satisfies Fact,
// 	// ImportObjectFact copies the value to *ptr.
// 	//
// 	// ImportObjectFact panics if called after the pass is complete.
// 	// ImportObjectFact is not concurrency-safe.
// 	ImportObjectFact func(obj types.Object, fact Fact) bool

// 	// ImportPackageFact retrieves a fact associated with package pkg,
// 	// which must be this package or one of its dependencies.
// 	// See comments for ImportObjectFact.
// 	ImportPackageFact func(pkg *types.Package, fact Fact) bool

// 	// ExportObjectFact associates a fact of type *T with the obj,
// 	// replacing any previous fact of that type.
// 	//
// 	// ExportObjectFact panics if it is called after the pass is
// 	// complete, or if obj does not belong to the package being analyzed.
// 	// ExportObjectFact is not concurrency-safe.
// 	ExportObjectFact func(obj types.Object, fact Fact)

// 	// ExportPackageFact associates a fact with the current package.
// 	// See comments for ExportObjectFact.
// 	ExportPackageFact func(fact Fact)

// 	// AllPackageFacts returns a new slice containing all package facts of the analysis's FactTypes
// 	// in unspecified order.
// 	// WARNING: This is an experimental API and may change in the future.
// 	AllPackageFacts func() []PackageFact

// 	// AllObjectFacts returns a new slice containing all object facts of the analysis's FactTypes
// 	// in unspecified order.
// 	// WARNING: This is an experimental API and may change in the future.
// 	AllObjectFacts func() []ObjectFact
// }

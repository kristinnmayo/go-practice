package main

import (
	"os"

	"github.com/go-practice/cmd"
	"github.com/go-practice/target"
)

func main() {
	// todo -- command line options
	path := os.Args[1]
	targets := target.Walk(path)

	// cmd.Tryscanner(targets)
	cmd.Astseek(targets)
	// cmd.Tryanalysis(targets)

	// log.Println("Targets:")
	// for _, t := range targets {
	// 	t.Display()
	// }
}

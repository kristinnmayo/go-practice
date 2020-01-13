package main

import (
	"github.com/go-practice/cmd"
)

func main() {
	// todo -- command line options
	// path := os.Args[1]
	// targets := target.Walk(path)

	// cmd.Scannerseek(targets)
	// cmd.Astseek(targets)
	cmd.Analysisseek()
	// cmd.Ssaseek(targets)
	// cmd.Customseek(targets)

	// log.Println("Targets:")
	// for _, t := range targets {
	// 	t.Display()
	// }
}

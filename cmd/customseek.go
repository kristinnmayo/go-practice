package cmd

import (
	"log"

	"github.com/go-practice/types/target"
)

// Customseek ...
func Customseek(targets []*target.Target) {
	for _, t := range targets {
		log.Println(t)
	}
}

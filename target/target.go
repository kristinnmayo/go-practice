package target

import (
	"log"
	"os"
)

// Target is one file in the user specified system
type Target struct {
	Path  string         // file path
	Perm  os.FileMode    // file permissions
	Vulns map[string]int // number of hits found
}

// New is used to create a new target
func New(path string, perm os.FileMode) Target {
	return Target{path, perm, make(map[string]int)}
}

// Display is used to log info on a target
func (t *Target) Display() {
	log.Printf("%s (%04o)", t.Path, t.Perm)
	log.Println(t.Vulns)
}

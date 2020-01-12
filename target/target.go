package target

import (
	"log"
	"os"
	"path/filepath"
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

// Walk accepts a path and return all filenames in tree
func Walk(dir string) (targets []*Target) {
	// func Walk(root string, walkFn WalkFunc) error
	// type WalkFunc func(path string, info os.FileInfo, err error) error
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatal(path, err)
		}
		// only check .go files
		if filepath.Ext(path) != ".go" {
			return nil
		}
		t := New(path, info.Mode().Perm())
		targets = append(targets, &t)
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	return targets
}

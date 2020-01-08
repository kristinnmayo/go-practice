package hit

import (
	"log"
)

// Hit is a possible vulnerability
type Hit struct {
	File string // path to file
	Line int    // line number
	Code string // vulnerable line of code
	Vuln string // type of vulnerability
}

// New is used to create a new hit
func New(file, code, vuln string, line int) Hit {
	return Hit{file, line, code, vuln}
}

// Display is used to log info on a hit
func (h Hit) Display() {
	log.Printf("%s:%d %s (%s)", h.File, h.Line, h.Code, h.Vuln)
}

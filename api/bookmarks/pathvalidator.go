package bookmarks

import (
	"fmt"
	"strings"
)

// Validator validates the existance of a folder with a given path and name
// e.g. /A_B_C/D  - path: /A_B_C, name: D
type Validator interface {
	Exists(path, name string) bool
}

// ValidatePath is used to check the existence of a path. The method uses
// recursion and starts with the 'last' path and traverses 'up' the tree
// e.g. path: /a/b/c/d
// validate, if a Folder node is present for
//	- segment: /a/b/c/d 	-> path: /a/b/c,	name: d
//	- segment: /a/b/c 	-> path: /a/b, 		name: c
//	- segment: /a/b 	-> path: /a, 		name: b
//	- segment: /a 		-> path: /, 		name: a
func ValidatePath(path string, check Validator) error {
	if path == "" {
		return fmt.Errorf("cannot use empty path")
	}
	if path == "/" {
		// the start of the whole path is valid without having a Folder node
		return nil
	}
	i := strings.LastIndex(path, "/")
	if i == -1 {
		return fmt.Errorf("no path seperator found")
	}
	p := path[0:i]
	n := path[i+1:]
	if i == 0 {
		p = "/"
	}
	e := check.Exists(p, n)
	if !e {
		return fmt.Errorf("the folder with path '%s' and name '%s' does not exist", p, n)
	}
	if i > 0 {
		return ValidatePath(path[0:i], check)
	}
	return nil
}

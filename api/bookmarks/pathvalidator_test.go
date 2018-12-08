package bookmarks_test

import (
	"testing"

	"github.com/bihe/bookmarks/api/bookmarks"
)

type simpleValidator struct{}

func (s simpleValidator) Exists(path, name string) bool {
	return true
}

func TestPathValidation(t *testing.T) {
	// Happy path
	tt := []struct {
		name string
		path string
	}{
		{
			name: "Multi-Path",
			path: "/a/b/c/d",
		},
		{
			name: "Longer-Path",
			path: "/A longer Path/to/check for the algorithm/how recursion works",
		},
		{
			name: "Singl-Path",
			path: "/a",
		},
		{
			name: "Root-Path",
			path: "/",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			err := bookmarks.ValidatePath(tc.path, simpleValidator{})
			if err != nil {
				t.Errorf("could not validate given path '%s': %v", tc.path, err)
			}
		})
	}

	// Error path
	tt = []struct {
		name string
		path string
	}{
		{
			name: "Empty",
			path: "",
		},
		{
			name: "Wrong",
			path: "\\a\\b\\d\\test path",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			err := bookmarks.ValidatePath(tc.path, simpleValidator{})
			if err == nil {
				t.Errorf("given path '%s' should be invalid!", tc.path)
			}
		})
	}

}

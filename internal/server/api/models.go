package api

import (
	"fmt"
	"net/http"
	"time"
)

// --------------------------------------------------------------------------
// Models
// --------------------------------------------------------------------------

type NodeType string

const (
	Node   NodeType = "Node"
	Folder NodeType = "Folder"
)

// Bookmark is the model provided via the REST API
// swagger:model
type Bookmark struct {
	ID          string     `json:"id"`
	Path        string     `json:"path"`
	DisplayName string     `json:"displayName"`
	URL         string     `json:"url"`
	SortOrder   int        `json:"sortOrder"`
	Type        NodeType   `json:"type"`
	Created     time.Time  `json:"created"`
	Modified    *time.Time `json:"modified,omitempty"`
	ChildCount  int        `json:"childCount"`
	AccessCount int        `json:"accessCount"`
	Favicon     string     `json:"favicon"`
}

// --------------------------------------------------------------------------
// BookmarkRequest
// --------------------------------------------------------------------------

// BookmarkRequest is the request payload for Bookmark data model.
type BookmarkRequest struct {
	*Bookmark
}

// Bind assigns the the provided data to a BookmarkRequest
func (b *BookmarkRequest) Bind(r *http.Request) error {
	// a Bookmark is nil if no Bookmark fields are sent in the request. Return an
	// error to avoid a nil pointer dereference.
	if b.Bookmark == nil {
		return fmt.Errorf("missing required Bookmarks fields")
	}
	return nil
}

// --------------------------------------------------------------------------
// BookmarkResponse
// --------------------------------------------------------------------------

// BookmarkResponse is the response payload for the Bookmark data model.
//
// In the BookmarkResponse object, first a Render() is called on itself,
// then the next field, and so on, all the way down the tree.
// Render is called in top-down order, like a http handler middleware chain.
type BookmarkResponse struct {
	*Bookmark
}

// Render the specific response
func (b BookmarkResponse) Render(w http.ResponseWriter, r *http.Request) error {
	// Pre-processing before a response is marshalled and sent across the wire
	return nil
}

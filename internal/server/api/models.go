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

// BookmarkList is a collection of Bookmarks
// swagger:model
type BookmarkList struct {
	Success bool       `json:"success"`
	Count   int        `json:"count"`
	Message string     `json:"message"`
	Value   []Bookmark `json:"value"`
}

// BookmarkResult hast additional information about a Bookmark
// swagger:model
type BookmarkResult struct {
	Success bool     `json:"success"`
	Message string   `json:"message"`
	Value   Bookmark `json:"value"`
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
type BookmarkResponse struct {
	*Bookmark
}

// Render the specific response
func (b BookmarkResponse) Render(w http.ResponseWriter, r *http.Request) error {
	// Pre-processing before a response is marshalled and sent across the wire
	return nil
}

// --------------------------------------------------------------------------
// BookmarkListResponse
// --------------------------------------------------------------------------

// BookmarkListResponse returns a list of Bookmark Items
type BookmarkListResponse struct {
	*BookmarkList
}

// Render the specific response
func (b BookmarkListResponse) Render(w http.ResponseWriter, r *http.Request) error {
	// Pre-processing before a response is marshalled and sent across the wire
	return nil
}

// --------------------------------------------------------------------------
// BookmarResultResponse
// --------------------------------------------------------------------------

// BookmarResultResponse returns BookmarResult
type BookmarResultResponse struct {
	*BookmarkResult
}

// Render the specific response
func (b BookmarResultResponse) Render(w http.ResponseWriter, r *http.Request) error {
	// Pre-processing before a response is marshalled and sent across the wire
	return nil
}

package bookmarks

import (
	"errors"
	"fmt"
	"net/http"
)

// --------------------------------------------------------------------------
// define API request / response structures
// --------------------------------------------------------------------------

const (
	// Node is a single bookmark entry
	Node string = "node"
	// Folder is a grouping/hierarchy structure to hold bookmarks
	Folder string = "folder"
)

// Bookmark is the view-model returned by the API
type Bookmark struct {
	Path        string `json:"path"`
	DisplayName string `json:"displayName"`
	URL         string `json:"url"`
	NodeID      string `json:"nodeId,omitempty"`
	SortOrder   uint8  `json:"sortOrder"`
	Type        string `json:"type"`
}

// Validate the bookmark object based on required fields
func (b Bookmark) Validate() error {
	if b.Path == "" {
		return fmt.Errorf("cannot use an empty path")
	}
	if b.Type == Node && b.URL == "" {
		return fmt.Errorf("a bookmarks needs an URL")
	}

	return nil
}

// BookmarkRequest is the request payload for Bookmark data model.
type BookmarkRequest struct {
	*Bookmark
}

// Bind assignes the the provided data to a BookmarkRequest
func (b *BookmarkRequest) Bind(r *http.Request) error {
	// a Bookmark is nil if no Bookmark fields are sent in the request. Return an
	// error to avoid a nil pointer dereference.
	if b.Bookmark == nil {
		return errors.New("missing required Bookmarks fields")
	}
	return nil
}

// BookmarkResponse is the response payload for the Bookmark data model.
//
// In the BookmarkResponse object, first a Render() is called on itself,
// then the next field, and so on, all the way down the tree.
// Render is called in top-down order, like a http handler middleware chain.
type BookmarkResponse struct {
	*Bookmark
}

// NewBookmarkResponse creates the response object needed for render
func NewBookmarkResponse(bookmark Bookmark) BookmarkResponse {
	resp := BookmarkResponse{Bookmark: &bookmark}
	return resp
}

// Render the specific response
func (b BookmarkResponse) Render(w http.ResponseWriter, r *http.Request) error {
	// Pre-processing before a response is marshalled and sent across the wire
	return nil
}

// BookmarkListResponse defines a list type
type BookmarkListResponse struct {
	Count int                `json:"count"`
	List  []BookmarkResponse `json:"result"`
}

// Render the specific response
func (b BookmarkListResponse) Render(w http.ResponseWriter, r *http.Request) error {
	// Pre-processing before a response is marshalled and sent across the wire
	return nil
}

// NewBookmarkListResponse creates the response for a list of objects
func NewBookmarkListResponse(bookmarks []Bookmark) BookmarkListResponse {
	var list = make([]BookmarkResponse, 0)
	for _, bookmark := range bookmarks {
		list = append(list, NewBookmarkResponse(bookmark))
	}
	resp := BookmarkListResponse{Count: len(list), List: list}
	return resp
}

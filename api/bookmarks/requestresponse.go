package bookmarks

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
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
	Created     int32  `json:"created"`
	Modified    int32  `json:"modified"`
	UserName    string `json:"username"`
}

var invalCharsDisplayName = []string{"/", "?", "\\", "\"", "<", ">", "#", "%", "{", "}", "|", "\\", "^", "~", "`", ";", "@", ":", "=", "&"}

// Validate the bookmark object based on required fields
func (b Bookmark) Validate() error {
	if b.Path == "" {
		return fmt.Errorf("cannot use an empty path")
	}
	if strings.HasSuffix(b.Path, "/") && b.Path != "/" {
		return fmt.Errorf("a path cannot end with '/")
	}
	if b.Type == Node && b.URL == "" {
		return fmt.Errorf("a bookmarks needs an URL")
	}
	for _, c := range invalCharsDisplayName {
		if strings.ContainsAny(b.DisplayName, c) {
			return fmt.Errorf("invalid chars in 'DisplayName'")
		}
	}
	return nil
}

// --------------------------------------------------------------------------
// BookmarkRequest
// --------------------------------------------------------------------------

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

// --------------------------------------------------------------------------
// BookmarkCreatedResponse
// --------------------------------------------------------------------------

// BookmarkCreatedResponse defines a generic success status
type BookmarkCreatedResponse struct {
	HTTPStatusCode int    `json:"status"`            // http response status code
	Message        string `json:"message,omitempty"` // application-level message
	NodeID         string `json:"nodeId,omitempty"`  // the ID of the creted bookmark
}

// NewBookmarkCreatedResponse created a success result
func NewBookmarkCreatedResponse(code int, message, nodeID string) *BookmarkCreatedResponse {
	return &BookmarkCreatedResponse{
		HTTPStatusCode: code,
		Message:        message,
		NodeID:         nodeID,
	}
}

// Render is the overloaded method for the ErrResponse
func (s *BookmarkCreatedResponse) Render(w http.ResponseWriter, r *http.Request) error {
	// Pre-processing before a response is marshalled and sent across the wire
	return nil
}

// --------------------------------------------------------------------------
// BookmarkListResponse
// --------------------------------------------------------------------------

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

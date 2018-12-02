package models

import (
	"errors"
	"net/http"

	"github.com/go-chi/render"
)

// Bookmark is the view-model returned by the API
type Bookmark struct {
	Path        string `json:"path"`
	DisplayName string `json:"displayName"`
	URL         string `json:"url"`
	NodeID      string `json:"nodeId,omitempty"`
	SortOrder   uint8  `json:"sortOrder"`
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
func NewBookmarkResponse(bookmark *Bookmark) *BookmarkResponse {
	resp := &BookmarkResponse{Bookmark: bookmark}
	return resp
}

// Render the specific response
func (b *BookmarkResponse) Render(w http.ResponseWriter, r *http.Request) error {
	// Pre-processing before a response is marshalled and sent across the wire
	return nil
}

// BookmarkListResponse defines a list type
type BookmarkListResponse struct {
	Count int                 `json:"count"`
	List  []*BookmarkResponse `json:"result"`
}

// Render the specific response
func (b *BookmarkListResponse) Render(w http.ResponseWriter, r *http.Request) error {
	// Pre-processing before a response is marshalled and sent across the wire
	return nil
}

// NewBookmarkListResponse creates the response for a list of objects
func NewBookmarkListResponse(bookmarks []*Bookmark) *BookmarkListResponse {
	var list []*BookmarkResponse
	for _, bookmark := range bookmarks {
		list = append(list, NewBookmarkResponse(bookmark))
	}
	resp := &BookmarkListResponse{Count: len(list), List: list}
	return resp
}

// SuccessResponse defines a generic success status
type SuccessResponse struct {
	HTTPStatusCode int    `json:"status"`            // http response status code
	Message        string `json:"message,omitempty"` // application-level message
}

// SuccessResult created a success result
func SuccessResult(code int, message string) render.Renderer {
	return &SuccessResponse{
		HTTPStatusCode: code,
		Message:        message,
	}
}

// Render is the overloaded method for the ErrResponse
func (s *SuccessResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, s.HTTPStatusCode)
	return nil
}

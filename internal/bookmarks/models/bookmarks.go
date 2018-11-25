package models

import (
	"net/http"
)

// Bookmark is the view-model returned by the API
type Bookmark struct {
	Path        string `json:"path"`
	DisplayName string `json:"displayName"`
	URL         string `json:"url"`
	NodeID      string `json:"nodeId"`
	SortOrder   uint8  `json:"sortOrder"`
	ItemType    string `json:"itemType"`
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
func (rd *BookmarkResponse) Render(w http.ResponseWriter, r *http.Request) error {
	// Pre-processing before a response is marshalled and sent across the wire
	return nil
}

// BookmarkListResponse defines a list type
type BookmarkListResponse struct {
	Count int                 `json:"count"`
	List  []*BookmarkResponse `json:"result"`
}

// Render the specific response
func (rd *BookmarkListResponse) Render(w http.ResponseWriter, r *http.Request) error {
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

package bookmarks

import (
	"fmt"
	"net/http"

	"github.com/bihe/bookmarks/api"
	"github.com/bihe/bookmarks/core"
	"github.com/bihe/bookmarks/security"
	"github.com/bihe/bookmarks/store"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

// --------------------------------------------------------------------------
// Bookmark API
// --------------------------------------------------------------------------

// BookmarkAPI combines the API methods of the bookmarks logic
type BookmarkAPI struct {
	uow *store.UnitOfWork
}

// GetAll retrieves the complete list of bookmarks entries from the store
func (app *BookmarkAPI) GetAll(w http.ResponseWriter, r *http.Request) {
	var err error
	var bookmarks = make([]store.BookmarkItem, 0)
	if bookmarks, err = app.uow.AllBookmarks(); err != nil {
		render.Render(w, r, api.ErrNotFound(err))
		return
	}
	render.Render(w, r, NewBookmarkListResponse(mapBookmarks(bookmarks)))
}

// GetByID returns a single bookmark item, path param :NodeId
func (app *BookmarkAPI) GetByID(w http.ResponseWriter, r *http.Request) {
	nodeID := chi.URLParam(r, "NodeID")
	if nodeID == "" {
		render.Render(w, r, api.ErrInvalidRequest(fmt.Errorf("missing id to load bookmark")))
		return
	}
	var item *store.BookmarkItem
	var err error
	if item, err = app.uow.BookmarkByID(nodeID); err != nil {
		render.Render(w, r, api.ErrNotFound(err))
		return
	}
	render.Render(w, r, NewBookmarkResponse(mapBookmark(*item)))
}

// FindByPath returns bookmarks/folders with the given path
// the path to find is provided as a query string
func (app *BookmarkAPI) FindByPath(w http.ResponseWriter, r *http.Request) {
	var err error
	var bookmarks []store.BookmarkItem
	path := r.URL.Query().Get("path")
	if path == "" {
		render.Render(w, r, api.ErrInvalidRequest(fmt.Errorf("no path supplied or missing query-param 'path'")))
		return
	}
	if bookmarks, err = app.uow.BookmarkByPath(path); err != nil {
		render.Render(w, r, api.ErrNotFound(err))
		return
	}
	render.Render(w, r, NewBookmarkListResponse(mapBookmarks(bookmarks)))
}

// Create will save a new bookmark entry
func (app *BookmarkAPI) Create(w http.ResponseWriter, r *http.Request) {
	var bookmark *Bookmark
	data := &BookmarkRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, api.ErrInvalidRequest(err))
		return
	}
	bookmark = data.Bookmark
	if err := bookmark.Validate(); err != nil {
		render.Render(w, r, api.ErrInvalidRequest(fmt.Errorf("invalid bookmark object: %v", err)))
		return
	}
	if bookmark.DisplayName == "" {
		bookmark.DisplayName = bookmark.URL
	}
	t := store.Node
	switch bookmark.Type {
	case Node:
		t = store.Node
	case Folder:
		t = store.Folder
	}

	err := app.uow.CreateBookmark(store.BookmarkItem{
		DisplayName: bookmark.DisplayName,
		Path:        bookmark.Path,
		URL:         bookmark.URL,
		SortOrder:   bookmark.SortOrder,
		Type:        t,
	})
	if err != nil {
		render.Render(w, r, api.ErrInvalidRequest(err))
		return
	}
	render.Render(w, r, api.SuccessResult(http.StatusCreated, fmt.Sprintf("bookmark item created: %s/%s", bookmark.Path, bookmark.DisplayName)))
}

// Update a bookmark item with new values
func (app *BookmarkAPI) Update(w http.ResponseWriter, r *http.Request) {
	var bookmark *Bookmark
	data := &BookmarkRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, api.ErrInvalidRequest(err))
		return
	}
	if data.Bookmark.NodeID == "" {
		render.Render(w, r, api.ErrInvalidRequest(fmt.Errorf("cannot upate bookmark with empty ID")))
		return
	}
	bookmark = data.Bookmark
	if err := bookmark.Validate(); err != nil {
		render.Render(w, r, api.ErrInvalidRequest(fmt.Errorf("invalid bookmark object: %v", err)))
		return
	}
	if bookmark.DisplayName == "" {
		bookmark.DisplayName = bookmark.URL
	}

	err := app.uow.UpdateBookmark(store.BookmarkItem{
		ItemID:      bookmark.NodeID,
		DisplayName: bookmark.DisplayName,
		Path:        bookmark.Path,
		URL:         bookmark.URL,
		SortOrder:   bookmark.SortOrder,
	})
	if err != nil {
		render.Render(w, r, api.ErrInvalidRequest(err))
		return
	}
	render.Render(w, r, api.SuccessResult(http.StatusOK, fmt.Sprintf("bookmark item updated: %s/%s", bookmark.Path, bookmark.DisplayName)))
}

// --------------------------------------------------------------------------
// internal helpers
// --------------------------------------------------------------------------

func mapBookmark(item store.BookmarkItem) Bookmark {
	var t string
	switch item.Type {
	case store.Node:
		t = Node
	case store.Folder:
		t = Folder
	}
	return Bookmark{
		DisplayName: item.DisplayName,
		Path:        item.Path,
		NodeID:      item.ItemID,
		SortOrder:   item.SortOrder,
		URL:         item.URL,
		Type:        t,
	}
}

func mapBookmarks(vs []store.BookmarkItem) []Bookmark {
	vsm := make([]Bookmark, len(vs))
	for i, v := range vs {
		vsm[i] = mapBookmark(v)
	}
	return vsm
}

func user(r *http.Request) *security.User {
	user := r.Context().Value(core.ContextUser).(*security.User)
	if user == nil {
		panic("could not get User from context")
	}
	return user
}

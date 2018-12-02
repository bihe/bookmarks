package bookmarks

import (
	"fmt"
	"net/http"

	"github.com/bihe/bookmarks-go/bookmarks/conf"
	"github.com/bihe/bookmarks-go/bookmarks/models"
	"github.com/bihe/bookmarks-go/bookmarks/security"
	"github.com/bihe/bookmarks-go/bookmarks/store"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

// BookmarkController combines the API methods of the bookmarks logic
type BookmarkController struct {
	uow *store.UnitOfWork
}

// GetAll retrieves the complete list of bookmarks entries from the store
func (app *BookmarkController) GetAll(w http.ResponseWriter, r *http.Request) {
	var err error
	var bookmarks []store.BookmarkItem
	if bookmarks, err = app.uow.GetAllBookmarks(); err != nil {
		render.Render(w, r, models.ErrNotFound(err))
	}
	render.Render(w, r, models.NewBookmarkListResponse(mapBookmarks(bookmarks)))
}

// Create will save a new bookmark entry
func (app *BookmarkController) Create(w http.ResponseWriter, r *http.Request) {
	var bookmark *models.Bookmark
	data := &models.BookmarkRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, models.ErrInvalidRequest(err))
		return
	}
	bookmark = data.Bookmark

	err := app.uow.CreateBookmark(store.BookmarkItem{
		DisplayName: bookmark.DisplayName,
		Path:        bookmark.Path,
		URL:         bookmark.URL,
		SortOrder:   bookmark.SortOrder,
	})
	if err != nil {
		render.Render(w, r, models.ErrInvalidRequest(err))
		return
	}
	render.Render(w, r, models.SuccessResult(http.StatusCreated, fmt.Sprintf("bookmark item created: %s/%s", bookmark.Path, bookmark.DisplayName)))
}

// Update a bookmark item with new values
func (app *BookmarkController) Update(w http.ResponseWriter, r *http.Request) {
	var bookmark *models.Bookmark
	data := &models.BookmarkRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, models.ErrInvalidRequest(err))
		return
	}
	if data.Bookmark.NodeID == "" {
		render.Render(w, r, models.ErrInvalidRequest(fmt.Errorf("cannot upate bookmark with empty ID")))
		return
	}
	bookmark = data.Bookmark

	err := app.uow.UpdateBookmark(store.BookmarkItem{
		ItemID:      bookmark.NodeID,
		DisplayName: bookmark.DisplayName,
		Path:        bookmark.Path,
		URL:         bookmark.URL,
		SortOrder:   bookmark.SortOrder,
	})
	if err != nil {
		render.Render(w, r, models.ErrInvalidRequest(err))
		return
	}
	render.Render(w, r, models.SuccessResult(http.StatusOK, fmt.Sprintf("bookmark item updated: %s/%s", bookmark.Path, bookmark.DisplayName)))
}

// GetByID returns a single bookmark item, path param :NodeId
func (app *BookmarkController) GetByID(w http.ResponseWriter, r *http.Request) {
	nodeID := chi.URLParam(r, "NodeID")
	if nodeID == "" {
		render.Render(w, r, models.ErrInvalidRequest(fmt.Errorf("missing id to load bookmark")))
		return
	}
	var item *store.BookmarkItem
	var err error
	if item, err = app.uow.GetItemByID(nodeID); err != nil {
		render.Render(w, r, models.ErrNotFound(err))
		return
	}
	render.Render(w, r, models.NewBookmarkResponse(mapBookmark(item)))
}

func mapBookmark(item *store.BookmarkItem) *models.Bookmark {
	return &models.Bookmark{
		DisplayName: item.DisplayName,
		Path:        item.Path,
		NodeID:      item.ItemID,
		SortOrder:   item.SortOrder,
		URL:         item.URL,
	}
}

func mapBookmarks(vs []store.BookmarkItem) []*models.Bookmark {
	vsm := make([]*models.Bookmark, len(vs))
	for i, v := range vs {
		vsm[i] = mapBookmark(&v)
	}
	return vsm
}

func user(r *http.Request) *security.User {
	user := r.Context().Value(conf.ContextUser).(*security.User)
	if user == nil {
		panic("could not get User from context")
	}
	return user
}

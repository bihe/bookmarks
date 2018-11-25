package bookmarks

import (
	"fmt"
	"net/http"

	"github.com/bihe/bookmarks-go/internal/bookmarks/models"
	"github.com/bihe/bookmarks-go/internal/conf"
	"github.com/bihe/bookmarks-go/internal/security"
	"github.com/bihe/bookmarks-go/internal/store"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

// BookmarkController combines the API methods of the bookmarks logic
type BookmarkController struct{}

// GetAll retrieves the complete list of bookmarks entries from the store
func (app *BookmarkController) GetAll(w http.ResponseWriter, r *http.Request) {
	var err error
	var bookmarks []store.BookmarkItem
	if bookmarks, err = uow(r).GetAllBookmarks(); err != nil {
		render.Render(w, r, models.ErrNotFound)
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
	itemType := store.BookmarkNode
	if bookmark.ItemType == "folder" {
		itemType = store.BookmarkNode
	}

	err := uow(r).CreateBookmark(store.BookmarkItem{
		DisplayName: bookmark.DisplayName,
		Path:        bookmark.Path,
		Type:        itemType,
		URL:         bookmark.URL,
		SortOrder:   bookmark.SortOrder,
	})
	if err != nil {
		render.Render(w, r, models.ErrInvalidRequest(err))
		return
	}
	render.Render(w, r, models.SuccessResult(http.StatusCreated, fmt.Sprintf("bookmark item created: %s/%s", bookmark.Path, bookmark.DisplayName)))
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
	if item, err = uow(r).GetItemByID(nodeID); err != nil {
		render.Render(w, r, models.ErrInvalidRequest(err))
		return
	}
	render.Render(w, r, models.NewBookmarkResponse(mapBookmark(item)))
}

func mapBookmark(item *store.BookmarkItem) *models.Bookmark {
	t := ""
	switch item.Type {
	case store.BookmarkFolder:
		t = "folder"
	default:
		t = "node"
	}
	return &models.Bookmark{
		DisplayName: item.DisplayName,
		Path:        item.Path,
		NodeID:      item.ItemID,
		ItemType:    t,
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

func uow(r *http.Request) *store.UnitOfWork {
	uow := r.Context().Value(conf.ContextUnitOfWork).(*store.UnitOfWork)
	if uow == nil {
		panic("could not get UnitOfWork from context")
	}
	return uow
}

func user(r *http.Request) *security.User {
	user := r.Context().Value(conf.ContextUser).(*security.User)
	if user == nil {
		panic("could not get User from context")
	}
	return user
}

package bookmarks

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/bihe/bookmarks/api/context"
	"github.com/bihe/bookmarks/api/models"
	"github.com/bihe/bookmarks/store"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/microcosm-cc/bluemonday"
)

// --------------------------------------------------------------------------
// Bookmark API Routes
// --------------------------------------------------------------------------

// MountRoutes defines the application specific routes
func MountRoutes(store *store.UnitOfWork) http.Handler {
	b := newAPI(store)

	r := chi.NewRouter()
	r.Options("/", b.GetAllOptions)
	r.Get("/", b.GetAll)
	r.Post("/", b.Create)
	r.Put("/", b.Update)
	r.Get("/{NodeID}", b.GetByID)
	r.Get("/path", b.FindByPath)
	r.Delete("/{NodeID}", b.Delete)
	r.Delete("/{NodeID}/{Force}", b.Delete)
	r.Get("/search", b.FindByName)

	return r
}

// --------------------------------------------------------------------------
// Bookmark API
// --------------------------------------------------------------------------

type bookmarksAPI interface {
	GetByID(http.ResponseWriter, *http.Request)

	GetAllOptions(http.ResponseWriter, *http.Request)
	GetAll(http.ResponseWriter, *http.Request)
	FindByPath(http.ResponseWriter, *http.Request)
	FindByName(http.ResponseWriter, *http.Request)

	Create(http.ResponseWriter, *http.Request)
	Update(http.ResponseWriter, *http.Request)
	Delete(http.ResponseWriter, *http.Request)
}

func newAPI(store *store.UnitOfWork) bookmarksAPI {
	return &bAPI{uow: store, policy: bluemonday.UGCPolicy()}
}

// --------------------------------------------------------------------------
// Bookmark API Implementation
// --------------------------------------------------------------------------

type bAPI struct {
	uow *store.UnitOfWork
	// policy is used to sanitize user input - prevent XSS
	policy *bluemonday.Policy
}

type result struct {
	w http.ResponseWriter
	r *http.Request
}

func (r *result) single(b *store.BookmarkItem, err error) {
	if err != nil {
		render.Render(r.w, r.r, models.ErrNotFound(models.NotFoundError{Request: r.r, Err: err}))
		return
	}
	render.Render(r.w, r.r, NewBookmarkResponse(mapBookmark(*b)))
}

func (r *result) list(bs []store.BookmarkItem, err error) {
	if err != nil {
		render.Render(r.w, r.r, models.ErrNotFound(models.NotFoundError{Request: r.r, Err: err}))
		return
	}
	render.Render(r.w, r.r, NewBookmarkListResponse(mapBookmarks(bs)))
}

// GetByID returns a single bookmark item, path param :NodeId
func (app *bAPI) GetByID(w http.ResponseWriter, r *http.Request) {
	rs := &result{w: w, r: r}

	nodeID := chi.URLParam(r, "NodeID")
	if nodeID == "" {
		render.Render(w, r, models.ErrBadRequest(models.BadRequestError{Request: r, Err: fmt.Errorf("missing id to load bookmark")}))
		return
	}
	rs.single(app.uow.BookmarkByID(nodeID, context.User(r).Username))
}

// --------------------------------------------------------------------------
// return a list of bookmarks
// --------------------------------------------------------------------------

// GetAllOptions used for OPTIONS request
func (app *bAPI) GetAllOptions(w http.ResponseWriter, r *http.Request) {
	var err error
	if _, err = app.uow.AllBookmarks(context.User(r).Username); err != nil {
		render.Render(w, r, models.ErrNotFound(models.NotFoundError{Request: r, Err: err}))
		return
	}
}

// GetAll retrieves the complete list of bookmarks entries from the store
func (app *bAPI) GetAll(w http.ResponseWriter, r *http.Request) {
	rs := &result{w: w, r: r}
	var bookmarks = make([]store.BookmarkItem, 0)
	var err error
	bookmarks, err = app.uow.AllBookmarks(context.User(r).Username)
	rs.list(bookmarks, err)
}

// FindByPath returns bookmarks/folders with the given path
// the path to find is provided as a query string
func (app *bAPI) FindByPath(w http.ResponseWriter, r *http.Request) {
	var err error
	var bookmarks []store.BookmarkItem
	path := r.URL.Query().Get("~p")
	if path == "" {
		render.Render(w, r, models.ErrBadRequest(models.BadRequestError{Request: r, Err: fmt.Errorf("no path supplied or missing query-param 'path'")}))
		return
	}
	path = sanitizeInput(app.policy, path)
	if bookmarks, err = app.uow.BookmarkByPath(path, context.User(r).Username); err != nil {
		render.Render(w, r, models.ErrNotFound(models.NotFoundError{Request: r, Err: err}))
		return
	}
	render.Render(w, r, NewBookmarkListResponse(mapBookmarks(bookmarks)))
}

// FindByName search the bookmark items for an element with the given name
func (app *bAPI) FindByName(w http.ResponseWriter, r *http.Request) {
	var err error
	var bookmarks []store.BookmarkItem
	name := r.URL.Query().Get("~n")
	if name == "" {
		render.Render(w, r, models.ErrBadRequest(models.BadRequestError{Request: r, Err: fmt.Errorf("no name supplied or missing query-param 'name'")}))
		return
	}
	name = sanitizeInput(app.policy, name)
	if bookmarks, err = app.uow.BookmarkByName(name, context.User(r).Username); err != nil {
		render.Render(w, r, models.ErrNotFound(models.NotFoundError{Request: r, Err: err}))
		return
	}
	render.Render(w, r, NewBookmarkListResponse(mapBookmarks(bookmarks)))
}

// --------------------------------------------------------------------------
// change bookmarks
// --------------------------------------------------------------------------

// Create will save a new bookmark entry
// the bookmark entry has a given type. If the type is Folder, the URL is just ignored and not saved
// if a bookmark is created with a given path, the method first checks if the path is available.
// This availability check of the path determines if for each path segment a corresponding Folder node is
// available. If the node is missing, the Node or Folder cannot be created for this path
func (app *bAPI) Create(w http.ResponseWriter, r *http.Request) {
	var bookmark *Bookmark
	data := &BookmarkRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, models.ErrBadRequest(models.BadRequestError{Request: r, Err: err}))
		return
	}
	bookmark = data.Bookmark
	// validate the supplied bookmark data for mandatory fields, invalid chars, ...
	if err := bookmark.Validate(); err != nil {
		render.Render(w, r, models.ErrBadRequest(models.BadRequestError{Request: r, Err: err}))
		return
	}
	// check if the given folder-structure is available
	if err := ValidatePath(bookmark.Path, dbFolderValidator{uow: app.uow, user: context.User(r).Username}); err != nil {
		render.Render(w, r, models.ErrBadRequest(models.BadRequestError{
			Request: r,
			Err:     fmt.Errorf("cannot create item because of missing folder structure: %v", err)}))
		return
	}

	if bookmark.DisplayName == "" {
		bookmark.DisplayName = bookmark.URL
	}

	t := store.Node
	url := bookmark.URL
	switch bookmark.Type {
	case Node:
		t = store.Node
	case Folder:
		t = store.Folder
		url = ""
	}

	b, err := app.uow.CreateBookmark(store.BookmarkItem{
		DisplayName: sanitizeInput(app.policy, bookmark.DisplayName),
		Path:        sanitizeInput(app.policy, bookmark.Path),
		URL:         url,
		SortOrder:   bookmark.SortOrder,
		Type:        t,
		Username:    context.User(r).Username,
	})
	if err != nil {
		render.Render(w, r, models.ErrBadRequest(models.BadRequestError{Request: r, Err: err}))
		return
	}
	render.Render(w, r, NewBookmarkCreatedResponse(http.StatusCreated, fmt.Sprintf("bookmark item created: p:%s, n:%s", bookmark.Path, bookmark.DisplayName), b.ItemID))
}

// Update a bookmark item with new values. The type of the bookmark Node/Folder is not updated.
// It does not make any sense to change a bookmark Node with URL to a Folder
func (app *bAPI) Update(w http.ResponseWriter, r *http.Request) {
	var bookmark *Bookmark
	data := &BookmarkRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, models.ErrBadRequest(models.BadRequestError{Request: r, Err: err}))
		return
	}
	if data.Bookmark.NodeID == "" {
		render.Render(w, r, models.ErrBadRequest(models.BadRequestError{Request: r, Err: fmt.Errorf("cannot upate bookmark with empty ID")}))
		return
	}
	bookmark = data.Bookmark

	if _, err := app.uow.BookmarkByID(bookmark.NodeID, context.User(r).Username); err != nil {
		render.Render(w, r, models.ErrNotFound(models.NotFoundError{Request: r, Err: fmt.Errorf("bookmark with ID '%s' not available", bookmark.NodeID)}))
		return
	}
	// validate the supplied bookmark data for mandatory fields, invalid chars, ...
	if err := bookmark.Validate(); err != nil {
		render.Render(w, r, models.ErrBadRequest(models.BadRequestError{Request: r, Err: err}))
		return
	}
	// check if the given folder-structure is available
	if err := ValidatePath(bookmark.Path, dbFolderValidator{uow: app.uow, user: context.User(r).Username}); err != nil {
		render.Render(w, r, models.ErrBadRequest(models.BadRequestError{
			Request: r,
			Err:     fmt.Errorf("cannot update item because of missing folder structure: %v", err),
		}))
		return
	}
	if bookmark.DisplayName == "" {
		bookmark.DisplayName = bookmark.URL
	}
	// the URL for a Folder does not make any sense
	url := bookmark.URL
	if bookmark.Type == Folder {
		url = ""
	}
	err := app.uow.UpdateBookmark(store.BookmarkItem{
		ItemID:      bookmark.NodeID,
		DisplayName: sanitizeInput(app.policy, bookmark.DisplayName),
		Path:        sanitizeInput(app.policy, bookmark.Path),
		URL:         url,
		SortOrder:   bookmark.SortOrder,
		Username:    context.User(r).Username,
	})
	if err != nil {
		render.Render(w, r, models.ErrBadRequest(models.BadRequestError{Request: r, Err: err}))
		return
	}
	render.Render(w, r, models.SuccessResult(http.StatusOK, fmt.Sprintf("bookmark item updated: %s/%s", bookmark.Path, bookmark.DisplayName)))
}

// Delete removes a given bookmark. It is necessary to provide a specific
// item-ID as a path param :NodeId
// If the item is a folder, and there are more items available 'within' this folder
// an error is returned, indicating this situation
// To forcefully delete the item nevertheless, the optional path param :Force can be applied (bool)
func (app *bAPI) Delete(w http.ResponseWriter, r *http.Request) {
	nodeID := chi.URLParam(r, "NodeID")
	if nodeID == "" {
		render.Render(w, r, models.ErrBadRequest(models.BadRequestError{Request: r, Err: fmt.Errorf("missing id, cannot delete bookmark")}))
		return
	}

	var item *store.BookmarkItem
	var err error
	if item, err = app.uow.BookmarkByID(nodeID, context.User(r).Username); err != nil {
		render.Render(w, r, models.ErrNotFound(models.NotFoundError{Request: r, Err: err}))
		return
	}
	if item.Type == store.Node {
		err := app.uow.Delete(nodeID, context.User(r).Username)
		if err != nil {
			render.Render(w, r, models.ErrBadRequest(models.BadRequestError{Request: r, Err: err}))
			return
		}
		render.Render(w, r, models.SuccessResult(http.StatusOK, fmt.Sprintf("bookmark item was deleted: '%s'", nodeID)))
		return
	}

	// for folders it needs to be checked, if there are 'children' within this folder
	path := item.Path + "/" + item.DisplayName
	if strings.HasSuffix(item.Path, "/") {
		path = item.Path + item.DisplayName
	}

	var bookmarks []store.BookmarkItem
	if bookmarks, err = app.uow.BookmarkStartsByPath(path, context.User(r).Username); err != nil {
		render.Render(w, r, models.ErrServerError(models.ServerError{Request: r, Err: err}))
		return
	}
	if len(bookmarks) == 0 {
		// all is good, we can just delete the item
		err := app.uow.Delete(nodeID, context.User(r).Username)
		if err != nil {
			render.Render(w, r, models.ErrBadRequest(models.BadRequestError{Request: r, Err: err}))
			return
		}
		render.Render(w, r, models.SuccessResult(http.StatusOK, fmt.Sprintf("bookmark item was deleted: '%s'", nodeID)))
		return
	}
	force, _ := strconv.ParseBool(chi.URLParam(r, "Force"))
	if !force {
		render.Render(w, r, models.ErrBadRequest(models.BadRequestError{
			Request: r,
			Err:     fmt.Errorf("cannot delete the item because of '%d' child elements", len(bookmarks))}))
		return
	}
	// force:True is supplied, so no matter what, we will delete the whole path
	err = app.uow.DeletePath(path, context.User(r).Username)
	if err != nil {
		render.Render(w, r, models.ErrBadRequest(models.BadRequestError{
			Request: r,
			Err:     fmt.Errorf("cannot delete the item by force: %v", err)}))
		return
	}
	render.Render(w, r, models.SuccessResult(http.StatusOK, fmt.Sprintf("bookmark item was deleted: '%s'", nodeID)))
}

// --------------------------------------------------------------------------
// validate a given path by checking if the folder-structure is avail in DB
// --------------------------------------------------------------------------

// dbFolderValidator checks for the existence of a 'Folder' item with the
// given 'name' in the given 'path'
type dbFolderValidator struct {
	uow  *store.UnitOfWork
	user string
}

func (d dbFolderValidator) Exists(path, name string) bool {
	_, err := d.uow.FolderByPathName(path, name, d.user)
	if err != nil {
		return false
	}
	return true

}

// --------------------------------------------------------------------------
// internal
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
		Modified:    item.Modified,
		Created:     item.Created,
		UserName:    item.Username,
	}
}

func mapBookmarks(vs []store.BookmarkItem) []Bookmark {
	vsm := make([]Bookmark, len(vs))
	for i, v := range vs {
		vsm[i] = mapBookmark(v)
	}
	return vsm
}

// sanitizeInput removes unwanted elements from the user-input
func sanitizeInput(policy *bluemonday.Policy, input string) string {
	return policy.Sanitize(input)
}

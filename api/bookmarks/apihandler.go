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
	r.Options("/", check(b.GetAllOptions))
	r.Get("/", list(b.GetAll))
	r.Post("/", created(b.Create))
	r.Put("/", ok(b.Update))
	r.Get("/{NodeID}", single(b.GetByID))
	r.Get("/path", list(b.FindByPath))
	r.Delete("/{NodeID}", ok(b.Delete))
	r.Delete("/{NodeID}/{Force}", ok(b.Delete))
	r.Get("/search", list(b.FindByName))

	return r
}

// --------------------------------------------------------------------------
// Bookmark API
// --------------------------------------------------------------------------

type bookmarksAPI interface {
	GetByID(http.ResponseWriter, *http.Request) (*Bookmark, error)

	GetAllOptions(http.ResponseWriter, *http.Request) error
	GetAll(http.ResponseWriter, *http.Request) ([]Bookmark, error)
	FindByPath(http.ResponseWriter, *http.Request) ([]Bookmark, error)
	FindByName(http.ResponseWriter, *http.Request) ([]Bookmark, error)

	Create(http.ResponseWriter, *http.Request) (string, string, error)
	Update(http.ResponseWriter, *http.Request) (string, error)
	Delete(http.ResponseWriter, *http.Request) (string, error)
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

// GetByID returns a single bookmark item, path param :NodeId
func (app *bAPI) GetByID(w http.ResponseWriter, r *http.Request) (*Bookmark, error) {
	nodeID := chi.URLParam(r, "NodeID")
	if nodeID == "" {
		return nil, models.BadRequestError{Request: r, Err: fmt.Errorf("missing id to load bookmark")}
	}
	var (
		bookmark *store.BookmarkItem
		err      error
	)
	if bookmark, err = app.uow.BookmarkByID(nodeID, context.User(r).Username); err != nil {
		return nil, models.NotFoundError{Request: r, Err: err}
	}
	b := mapBookmark(*bookmark)
	return &b, nil
}

// GetAllOptions used for OPTIONS request
func (app *bAPI) GetAllOptions(w http.ResponseWriter, r *http.Request) error {
	var err error
	if _, err = app.uow.AllBookmarks(context.User(r).Username); err != nil {
		return models.NotFoundError{Request: r, Err: err}
	}
	return nil
}

// GetAll retrieves the complete list of bookmarks entries from the store
func (app *bAPI) GetAll(w http.ResponseWriter, r *http.Request) ([]Bookmark, error) {
	var (
		bookmarks = make([]store.BookmarkItem, 0)
		err       error
	)
	if bookmarks, err = app.uow.AllBookmarks(context.User(r).Username); err != nil {
		return nil, models.NotFoundError{Request: r, Err: err}
	}
	return mapBookmarks(bookmarks), nil
}

// FindByPath returns bookmarks/folders with the given path
// the path to find is provided as a query string
func (app *bAPI) FindByPath(w http.ResponseWriter, r *http.Request) ([]Bookmark, error) {
	var (
		bookmarks = make([]store.BookmarkItem, 0)
		err       error
	)
	path := r.URL.Query().Get("~p")
	if path == "" {
		return nil, models.BadRequestError{Request: r, Err: fmt.Errorf("no path supplied or missing query-param 'path'")}
	}
	path = sanitizeInput(app.policy, path)
	if bookmarks, err = app.uow.BookmarkByPath(path, context.User(r).Username); err != nil {
		return nil, models.NotFoundError{Request: r, Err: err}
	}
	return mapBookmarks(bookmarks), nil
}

// FindByName search the bookmark items for an element with the given name
func (app *bAPI) FindByName(w http.ResponseWriter, r *http.Request) ([]Bookmark, error) {
	var (
		bookmarks = make([]store.BookmarkItem, 0)
		err       error
	)
	name := r.URL.Query().Get("~n")
	if name == "" {
		return nil, models.BadRequestError{Request: r, Err: fmt.Errorf("no name supplied or missing query-param 'name'")}
	}
	name = sanitizeInput(app.policy, name)
	if bookmarks, err = app.uow.BookmarkByName(name, context.User(r).Username); err != nil {
		return nil, models.NotFoundError{Request: r, Err: err}
	}
	return mapBookmarks(bookmarks), nil
}

// --------------------------------------------------------------------------
// change bookmarks
// --------------------------------------------------------------------------

// Create will save a new bookmark entry
// the bookmark entry has a given type. If the type is Folder, the URL is just ignored and not saved
// if a bookmark is created with a given path, the method first checks if the path is available.
// This availability check of the path determines if for each path segment a corresponding Folder node is
// available. If the node is missing, the Node or Folder cannot be created for this path
func (app *bAPI) Create(w http.ResponseWriter, r *http.Request) (string, string, error) {
	var bookmark *Bookmark
	data := &BookmarkRequest{}
	if err := render.Bind(r, data); err != nil {
		return "", "", models.BadRequestError{Request: r, Err: err}
	}
	bookmark = data.Bookmark
	// validate the supplied bookmark data for mandatory fields, invalid chars, ...
	if err := bookmark.Validate(); err != nil {
		return "", "", models.BadRequestError{Request: r, Err: err}
	}
	// check if the given folder-structure is available
	if err := ValidatePath(bookmark.Path, dbFolderValidator{uow: app.uow, user: context.User(r).Username}); err != nil {
		return "", "", models.BadRequestError{
			Request: r,
			Err:     fmt.Errorf("cannot create item because of missing folder structure: %v", err)}
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
		return "", "", models.BadRequestError{Request: r, Err: err}
	}
	return fmt.Sprintf("bookmark item created: p:%s, n:%s", bookmark.Path, bookmark.DisplayName), b.ItemID, nil
}

// Update a bookmark item with new values. The type of the bookmark Node/Folder is not updated.
// It does not make any sense to change a bookmark Node with URL to a Folder
func (app *bAPI) Update(w http.ResponseWriter, r *http.Request) (string, error) {
	var bookmark *Bookmark
	data := &BookmarkRequest{}
	if err := render.Bind(r, data); err != nil {
		return "", models.BadRequestError{Request: r, Err: err}
	}
	if data.Bookmark.NodeID == "" {
		return "", models.BadRequestError{Request: r, Err: fmt.Errorf("cannot upate bookmark with empty ID")}
	}
	bookmark = data.Bookmark

	if _, err := app.uow.BookmarkByID(bookmark.NodeID, context.User(r).Username); err != nil {
		return "", models.NotFoundError{Request: r, Err: fmt.Errorf("bookmark with ID '%s' not available", bookmark.NodeID)}
	}
	// validate the supplied bookmark data for mandatory fields, invalid chars, ...
	if err := bookmark.Validate(); err != nil {
		return "", models.BadRequestError{Request: r, Err: err}
	}
	// check if the given folder-structure is available
	if err := ValidatePath(bookmark.Path, dbFolderValidator{uow: app.uow, user: context.User(r).Username}); err != nil {
		return "", models.BadRequestError{
			Request: r,
			Err:     fmt.Errorf("cannot update item because of missing folder structure: %v", err),
		}
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
		return "", models.BadRequestError{Request: r, Err: err}
	}
	return fmt.Sprintf("bookmark item updated: %s/%s", bookmark.Path, bookmark.DisplayName), nil
}

// Delete removes a given bookmark. It is necessary to provide a specific
// item-ID as a path param :NodeId
// If the item is a folder, and there are more items available 'within' this folder
// an error is returned, indicating this situation
// To forcefully delete the item nevertheless, the optional path param :Force can be applied (bool)
func (app *bAPI) Delete(w http.ResponseWriter, r *http.Request) (string, error) {
	nodeID := chi.URLParam(r, "NodeID")
	if nodeID == "" {
		return "", models.BadRequestError{Request: r, Err: fmt.Errorf("missing id, cannot delete bookmark")}
	}

	var item *store.BookmarkItem
	var err error
	if item, err = app.uow.BookmarkByID(nodeID, context.User(r).Username); err != nil {
		return "", models.NotFoundError{Request: r, Err: err}
	}
	if item.Type == store.Node {
		err := app.uow.Delete(nodeID, context.User(r).Username)
		if err != nil {
			return "", models.BadRequestError{Request: r, Err: err}
		}
		return fmt.Sprintf("bookmark item was deleted: '%s'", nodeID), nil
	}

	// for folders it needs to be checked, if there are 'children' within this folder
	path := item.Path + "/" + item.DisplayName
	if strings.HasSuffix(item.Path, "/") {
		path = item.Path + item.DisplayName
	}

	var bookmarks []store.BookmarkItem
	if bookmarks, err = app.uow.BookmarkStartsByPath(path, context.User(r).Username); err != nil {
		return "", models.ServerError{Request: r, Err: err}
	}
	if len(bookmarks) == 0 {
		// all is good, we can just delete the item
		err := app.uow.Delete(nodeID, context.User(r).Username)
		if err != nil {
			return "", models.BadRequestError{Request: r, Err: err}
		}
		return fmt.Sprintf("bookmark item was deleted: '%s'", nodeID), nil
	}
	force, _ := strconv.ParseBool(chi.URLParam(r, "Force"))
	if !force {
		return "", models.BadRequestError{
			Request: r,
			Err:     fmt.Errorf("cannot delete the item because of '%d' child elements", len(bookmarks))}
	}
	// force:True is supplied, so no matter what, we will delete the whole path
	err = app.uow.DeletePath(path, context.User(r).Username)
	if err != nil {
		return "", models.BadRequestError{
			Request: r,
			Err:     fmt.Errorf("cannot delete the item by force: %v", err)}
	}
	return fmt.Sprintf("bookmark item was deleted: '%s'", nodeID), nil
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

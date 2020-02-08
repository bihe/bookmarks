// Package api implements the HTTP API of the bookmarks-application.
//
// the purpose of this application is to provide bookmarks management independent of browsers
//
// Terms Of Service:
//
//     Schemes: https
//     Host: bookmarks.binggl.net
//     BasePath: /api/v1
//     Version: 1.0.0
//     License: Apache 2.0 https://opensource.org/licenses/Apache-2.0
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
// swagger:meta
package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/bihe/bookmarks/internal/store"
	"github.com/bihe/commons-go/errors"
	"github.com/bihe/commons-go/handler"
	"github.com/bihe/commons-go/security"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

// BookmarksAPI implement the handler for the API
type BookmarksAPI struct {
	handler.Handler
	Repository store.Repository
}

// swagger:operation GET /api/v1/bookmarks/{id} bookmarks GetBookmarkByID
//
// get a bookmark by id
//
// returns a single bookmark specified by it's ID
//
// ---
// produces:
// - application/json
// parameters:
// - name: id
//   in: path
// responses:
//   '200':
//     description: Bookmark
//     schema:
//       "$ref": "#/definitions/Bookmark"
//   '400':
//     description: ProblemDetail
//     schema:
//       "$ref": "#/definitions/ProblemDetail"
//   '404':
//     description: ProblemDetail
//     schema:
//       "$ref": "#/definitions/ProblemDetail"
//   '401':
//     description: ProblemDetail
//     schema:
//       "$ref": "#/definitions/ProblemDetail"
//   '403':
//     description: ProblemDetail
//     schema:
//       "$ref": "#/definitions/ProblemDetail"
func (b *BookmarksAPI) GetBookmarkByID(user security.User, w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	if id == "" {
		return errors.BadRequestError{Err: fmt.Errorf("missing id parameter"), Request: r}
	}

	handler.LogFunction("api.GetBookmarkByID").Debugf("try to get bookmark by ID: '%s' for user: '%s'", id, user.Username)

	bookmark, err := b.Repository.GetBookmarkById(id, user.Username)
	if err != nil {
		handler.LogFunction("api.GetBookmarkByID").Warnf("cannot get bookmark by ID: '%s', %v", id, err)
		return errors.NotFoundError{Err: fmt.Errorf("no bookmark with ID '%s' avaliable", id), Request: r}
	}

	return render.Render(w, r, BookmarkResponse{Bookmark: entityToModel(bookmark)})
}

// swagger:operation GET /api/v1/bookmarks/bypath bookmarks GetBookmarksByPath
//
// get bookmarks by path
//
// returns a list of bookmarks for a given path
//
// ---
// produces:
// - application/json
// parameters:
// - name: path
//   in: query
// responses:
//   '200':
//     description: BookmarkList
//     schema:
//       "$ref": "#/definitions/BookmarkList"
//   '400':
//     description: ProblemDetail
//     schema:
//       "$ref": "#/definitions/ProblemDetail"
//   '401':
//     description: ProblemDetail
//     schema:
//       "$ref": "#/definitions/ProblemDetail"
//   '403':
//     description: ProblemDetail
//     schema:
//       "$ref": "#/definitions/ProblemDetail"
func (b *BookmarksAPI) GetBookmarksByPath(user security.User, w http.ResponseWriter, r *http.Request) error {
	path := r.URL.Query().Get("path")

	if path == "" {
		return errors.BadRequestError{Err: fmt.Errorf("missing path parameter"), Request: r}
	}

	handler.LogFunction("api.GetBookmarksByPath").Debugf("get bookmarks by path: '%s' for user: '%s'", path, user.Username)

	bms, err := b.Repository.GetBookmarksByPath(path, user.Username)
	var bookmarks []Bookmark
	if err != nil {
		handler.LogFunction("api.GetBookmarksByPath").Warnf("cannot get bookmark by path: '%s', %v", path, err)
	}
	bookmarks = entityListToModel(bms)
	count := len(bookmarks)
	result := BookmarkList{
		Success: true,
		Count:   count,
		Message: fmt.Sprintf("Found %d items.", count),
		Value:   bookmarks,
	}

	return render.Render(w, r, BookmarkListResponse{BookmarkList: &result})
}

// swagger:operation GET /api/v1/bookmarks/foolder bookmarks GetBookmarksFolderByPath
//
// get bookmark folder by path
//
// returns the folder identified by the given path
//
// ---
// produces:
// - application/json
// parameters:
// - name: path
//   in: query
// responses:
//   '200':
//     description: BookmarkResult
//     schema:
//       "$ref": "#/definitions/BookmarkResult"
//   '400':
//     description: ProblemDetail
//     schema:
//       "$ref": "#/definitions/ProblemDetail"
//   '404':
//     description: ProblemDetail
//     schema:
//       "$ref": "#/definitions/ProblemDetail"
//   '401':
//     description: ProblemDetail
//     schema:
//       "$ref": "#/definitions/ProblemDetail"
//   '403':
//     description: ProblemDetail
//     schema:
//       "$ref": "#/definitions/ProblemDetail"
func (b *BookmarksAPI) GetBookmarksFolderByPath(user security.User, w http.ResponseWriter, r *http.Request) error {
	path := r.URL.Query().Get("path")

	if path == "" {
		return errors.BadRequestError{Err: fmt.Errorf("missing path parameter"), Request: r}
	}

	handler.LogFunction("api.GetBookmarksFolderByPath").Debugf("get bookmarks-folder by path: '%s' for user: '%s'", path, user.Username)

	if path == "/" {
		// special treatment for the root path. This path is ALWAYS available
		// and does not have a specific storage entry - this is by convention
		return render.Render(w, r, BookmarResultResponse{BookmarkResult: &BookmarkResult{
			Success: true,
			Message: fmt.Sprintf("Found bookmark folder for path %s.", path),
			Value: Bookmark{
				DisplayName: "Root",
				Path:        "/",
				Type:        Folder,
				ID:          fmt.Sprintf("%s_ROOT", user.Username),
			},
		}})
	}

	bm, err := b.Repository.GetFolderByPath(path, user.Username)
	if err != nil {
		handler.LogFunction("api.GetBookmarksFolderByPath").Warnf("cannot get bookmark folder by path: '%s', %v", path, err)
		return errors.NotFoundError{Err: fmt.Errorf("no folder for path '%s' found", path), Request: r}
	}

	return render.Render(w, r, BookmarResultResponse{BookmarkResult: &BookmarkResult{
		Success: true,
		Message: fmt.Sprintf("Found bookmark folder for path %s.", path),
		Value:   *entityToModel(bm),
	}})
}

// swagger:operation GET /api/v1/bookmarks/byname bookmarks GetBookmarksByName
//
// get bookmarks by name
//
// search for bookmarks by name and return a list of search-results
//
// ---
// produces:
// - application/json
// parameters:
// - name: name
//   in: query
// responses:
//   '200':
//     description: BookmarkList
//     schema:
//       "$ref": "#/definitions/BookmarkList"
//   '400':
//     description: ProblemDetail
//     schema:
//       "$ref": "#/definitions/ProblemDetail"
//   '401':
//     description: ProblemDetail
//     schema:
//       "$ref": "#/definitions/ProblemDetail"
//   '403':
//     description: ProblemDetail
//     schema:
//       "$ref": "#/definitions/ProblemDetail"
func (b *BookmarksAPI) GetBookmarksByName(user security.User, w http.ResponseWriter, r *http.Request) error {
	name := r.URL.Query().Get("name")

	if name == "" {
		return errors.BadRequestError{Err: fmt.Errorf("missing name parameter"), Request: r}
	}

	handler.LogFunction("api.GetBookmarksByName").Debugf("get bookmarks by name: '%s' for user: '%s'", name, user.Username)

	bms, err := b.Repository.GetBookmarksByName(name, user.Username)
	var bookmarks []Bookmark
	if err != nil {
		handler.LogFunction("api.GetBookmarksByName").Warnf("cannot get bookmark by name: '%s', %v", name, err)
	}
	bookmarks = entityListToModel(bms)
	count := len(bookmarks)
	result := BookmarkList{
		Success: true,
		Count:   count,
		Message: fmt.Sprintf("Found %d items.", count),
		Value:   bookmarks,
	}

	return render.Render(w, r, BookmarkListResponse{BookmarkList: &result})
}

// swagger:operation GET /api/v1/bookmarks/mostvisited/{num} bookmarks GetMostVisited
//
// get recent accessed bookmarks
//
// return the most recently visited bookmarks
//
// ---
// produces:
// - application/json
// parameters:
// - name: num
//   in: path
// responses:
//   '200':
//     description: BookmarkList
//     schema:
//       "$ref": "#/definitions/BookmarkList"
//   '401':
//     description: ProblemDetail
//     schema:
//       "$ref": "#/definitions/ProblemDetail"
//   '403':
//     description: ProblemDetail
//     schema:
//       "$ref": "#/definitions/ProblemDetail"
func (b *BookmarksAPI) GetMostVisited(user security.User, w http.ResponseWriter, r *http.Request) error {
	n := chi.URLParam(r, "num")
	num, _ := strconv.Atoi(n)
	if num < 1 {
		num = 100
	}

	handler.LogFunction("api.GetMostVisited").Debugf("get the most recent, most often visited bookmarks for user: '%s'", user.Username)

	bms, err := b.Repository.GetMostRecentBookmarks(user.Username, num)
	var bookmarks []Bookmark
	if err != nil {
		handler.LogFunction("api.GetMostVisited").Warnf("cannot get most visited bookmarks: '%v'", err)
	}
	bookmarks = entityListToModel(bms)
	count := len(bookmarks)
	result := BookmarkList{
		Success: true,
		Count:   count,
		Message: fmt.Sprintf("Found %d items.", count),
		Value:   bookmarks,
	}

	return render.Render(w, r, BookmarkListResponse{BookmarkList: &result})
}

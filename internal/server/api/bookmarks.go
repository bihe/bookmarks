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

// swagger:operation GET /bookmarks bookmarks GetBookmarkByID
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

	handler.LogFunction("api.GetBookmarkByID").Debugf("try to get bookmark by ID: '%s' for user: '%s'", id, user.Username)

	bookmark, err := b.Repository.GetBookmarkById(id, user.Username)
	if err != nil {
		handler.LogFunction("api.GetBookmarkByID").Warnf("cannot get bookmark by ID: '%s'", id)
		return errors.NotFoundError{Err: fmt.Errorf("no bookmark with ID '%s' avaliable", id), Request: r}
	}

	return render.Render(w, r, BookmarkResponse{Bookmark: entityToModel(bookmark)})
}

// swagger:operation GET /bookmarks/bypath bookmarks GetBookmarksByPath
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
func (b *BookmarksAPI) GetBookmarksByPath(user security.User, w http.ResponseWriter, r *http.Request) error {
	path := r.URL.Query().Get("path")

	handler.LogFunction("api.GetBookmarksByPath").Debugf("get bookmarks by path: '%s' for user: '%s'", path, user.Username)

	bms, err := b.Repository.GetBookmarksByPath(path, user.Username)
	var bookmarks []Bookmark
	if err != nil {
		handler.LogFunction("api.GetBookmarksByPath").Warnf("cannot get bookmark by path: '%s'", path)
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

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
	ID := chi.URLParam(r, "id")

	handler.LogFunction("api.GetBookmarkByID").Debugf("try to get bookmark by ID: '%s' for user: '%s'", ID, user.Username)

	bookmark, err := b.Repository.GetBookmarkById(ID, user.Username)
	if err != nil {
		handler.LogFunction("api.GetBookmarkByID").Warnf("cannot get bookmark by ID: '%s'", ID)
		return errors.NotFoundError{Err: fmt.Errorf("no bookmark with ID '%s' avaliable", ID), Request: r}
	}

	return render.Render(w, r, BookmarkResponse{Bookmark: entityToModel(bookmark)})
}

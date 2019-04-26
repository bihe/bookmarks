package bookmarks

import (
	"errors"
	"net/http"

	"github.com/bihe/bookmarks/api/models"
	"github.com/go-chi/render"
)

// --------------------------------------------------------------------------
// BookmarkRequest
// --------------------------------------------------------------------------

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

// --------------------------------------------------------------------------
// BookmarkResponse
// --------------------------------------------------------------------------

// BookmarkResponse is the response payload for the Bookmark data model.
//
// In the BookmarkResponse object, first a Render() is called on itself,
// then the next field, and so on, all the way down the tree.
// Render is called in top-down order, like a http handler middleware chain.
type BookmarkResponse struct {
	*Bookmark
}

// NewBookmarkResponse creates the response object needed for render
func NewBookmarkResponse(bookmark Bookmark) BookmarkResponse {
	resp := BookmarkResponse{Bookmark: &bookmark}
	return resp
}

// Render the specific response
func (b BookmarkResponse) Render(w http.ResponseWriter, r *http.Request) error {
	// Pre-processing before a response is marshalled and sent across the wire
	return nil
}

// --------------------------------------------------------------------------
// BookmarkCreatedResponse
// --------------------------------------------------------------------------

// BookmarkCreatedResponse defines a generic success status
type BookmarkCreatedResponse struct {
	HTTPStatusCode int    `json:"status"`            // http response status code
	Message        string `json:"message,omitempty"` // application-level message
	NodeID         string `json:"nodeId,omitempty"`  // the ID of the created bookmark
}

// NewBookmarkCreatedResponse created a success result
func NewBookmarkCreatedResponse(code int, message, nodeID string) *BookmarkCreatedResponse {
	return &BookmarkCreatedResponse{
		HTTPStatusCode: code,
		Message:        message,
		NodeID:         nodeID,
	}
}

// Render is the overloaded method for the ErrResponse
func (s *BookmarkCreatedResponse) Render(w http.ResponseWriter, r *http.Request) error {
	// Pre-processing before a response is marshalled and sent across the wire
	render.Status(r, s.HTTPStatusCode)
	return nil
}

// --------------------------------------------------------------------------
// BookmarkListResponse
// --------------------------------------------------------------------------

// BookmarkListResponse defines a list type
type BookmarkListResponse struct {
	Count int                `json:"count"`
	List  []BookmarkResponse `json:"result"`
}

// Render the specific response
func (b BookmarkListResponse) Render(w http.ResponseWriter, r *http.Request) error {
	// Pre-processing before a response is marshalled and sent across the wire
	return nil
}

// NewBookmarkListResponse creates the response for a list of objects
func NewBookmarkListResponse(bookmarks []Bookmark) BookmarkListResponse {
	var list = make([]BookmarkResponse, 0)
	for _, bookmark := range bookmarks {
		list = append(list, NewBookmarkResponse(bookmark))
	}
	resp := BookmarkListResponse{Count: len(list), List: list}
	return resp
}

// --------------------------------------------------------------------------
// BookmarkTreeResponse
// --------------------------------------------------------------------------

// BookmarkTreeResponse is used to return the whole bookmar tree
type BookmarkTreeResponse struct {
	*BookmarkTree
}

// NewBookmarkTreeResponse creates the response object needed for render
func NewBookmarkTreeResponse(tree BookmarkTree) BookmarkTreeResponse {
	resp := BookmarkTreeResponse{BookmarkTree: &tree}
	return resp
}

// Render the specific response
func (b BookmarkTreeResponse) Render(w http.ResponseWriter, r *http.Request) error {
	// Pre-processing before a response is marshalled and sent across the wire
	return nil
}

// --------------------------------------------------------------------------
// handlers which take care of setting the correct result and/or error
// --------------------------------------------------------------------------

func list(f func(http.ResponseWriter, *http.Request) ([]Bookmark, error)) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			bs  []Bookmark
			err error
		)
		if bs, err = f(w, r); err != nil {
			handleError(w, r, err)
			return
		}
		render.Render(w, r, NewBookmarkListResponse(bs))
	})
}

func single(f func(http.ResponseWriter, *http.Request) (*Bookmark, error)) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			b   *Bookmark
			err error
		)
		if b, err = f(w, r); err != nil {
			handleError(w, r, err)
			return
		}
		render.Render(w, r, NewBookmarkResponse(*b))
	})
}

func tree(f func(http.ResponseWriter, *http.Request) (*BookmarkTree, error)) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			b   *BookmarkTree
			err error
		)
		if b, err = f(w, r); err != nil {
			handleError(w, r, err)
			return
		}
		render.Render(w, r, NewBookmarkTreeResponse(*b))
	})
}

func ok(f func(http.ResponseWriter, *http.Request) (string, error)) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			result string
			err    error
		)
		if result, err = f(w, r); err != nil {
			handleError(w, r, err)
			return
		}
		render.Render(w, r, models.SuccessResult(http.StatusOK, result))
	})
}

func created(f func(http.ResponseWriter, *http.Request) (string, string, error)) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			result string
			id     string
			err    error
		)
		if result, id, err = f(w, r); err != nil {
			handleError(w, r, err)
			return
		}
		render.Render(w, r, NewBookmarkCreatedResponse(http.StatusCreated, result, id))
	})
}

func check(f func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			handleError(w, r, err)
		}
	})
}

func handleError(w http.ResponseWriter, r *http.Request, err error) {
	if err != nil {
		switch err.(type) {
		case models.NotFoundError:
			render.Render(w, r, models.ErrNotFound(err.(models.NotFoundError)))
		case models.BadRequestError:
			render.Render(w, r, models.ErrBadRequest(err.(models.BadRequestError)))
		case models.ServerError:
			render.Render(w, r, models.ErrServerError(err.(models.ServerError)))
		default:
			render.Render(w, r, models.ErrNotFound(models.NotFoundError{Request: r, Err: err}))
		}
	}
}

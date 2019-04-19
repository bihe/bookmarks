package bookmarks

import (
	"net/http"

	"github.com/bihe/bookmarks/api/models"
	"github.com/go-chi/render"
)

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

// --------------------------------------------------------------------------
// handle function to take a http.Handler and return a http.HandlerFunc
// --------------------------------------------------------------------------

func handle(h http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// before
		h.ServeHTTP(w, r)
		// after
	})
}

package models

import (
	"fmt"
	"net/http"

	"github.com/go-chi/render"
)

//--
// Error response payloads & renderers
//--

// ErrResponse renderer type for handling all sorts of errors.
//
// In the best case scenario, the excellent github.com/pkg/errors package
// helps reveal information on the error, setting it on Err, and in the Render()
// method, using it to set the application-specific error code in AppCode.
type ErrResponse struct {
	HTTPStatusCode int    `json:"status"` // http response status code
	Message        string `json:"message"`
	AppCode        int64  `json:"code,omitempty"`  // application-specific error code
	ErrorText      string `json:"error,omitempty"` // application-level error message, for debugging
}

// Render is the overloaded method for the ErrResponse
func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

// ErrInvalidRequest returns a http.StatusBadRequest
func ErrInvalidRequest(err error) render.Renderer {
	return &ErrResponse{
		HTTPStatusCode: http.StatusBadRequest,
		Message:        fmt.Sprintf("invalid request: %s", err),
	}
}

// ErrNotFound returns a http.StatusNotFound
func ErrNotFound(err error) render.Renderer {
	return &ErrResponse{
		HTTPStatusCode: http.StatusNotFound,
		Message:        fmt.Sprintf("cannot find item: %v", err),
	}
}

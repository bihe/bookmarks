package api

import (
	"fmt"
	"net/http"

	"github.com/go-chi/render"
)

// --------------------------------------------------------------------------
// Error responses adhering to https://tools.ietf.org/html/rfc7807
// --------------------------------------------------------------------------

// ProblemDetail combines the fields defined in RFC7807
//
// "Note that both "type" and "instance" accept relative URIs; this means
// that they must be resolved relative to the document's base URI"
type ProblemDetail struct {
	// Type is a URI reference [RFC3986] that identifies the
	// problem type.  This specification encourages that, when
	// dereferenced, it provide human-readable documentation for the problem
	Type string `json:"type"`
	// Title is a short, human-readable summary of the problem type
	Title string `json:"title"`
	// Status is the HTTP status code
	Status int `json:"status"`
	// Detail is a human-readable explanation specific to this occurrence of the problem
	Detail string `json:"detail,omitempty"`
	// Instance is a URI reference that identifies the specific occurrence of the problem
	Instance string `json:"instance,omitempty"`
}

// Render is the overloaded method for the ProblemDetail
func (p *ProblemDetail) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, p.Status)
	return nil
}

// --------------------------------------------------------------------------
// Specific Errors
// --------------------------------------------------------------------------

// NotFoundError is used when a given object cannot be found
type NotFoundError struct {
	Err     error
	Request *http.Request
}

// Error implements the error interface
func (e NotFoundError) Error() string {
	return fmt.Sprintf("the object for request '%s' cannot be found: %v", e.Request.RequestURI, e.Err)
}

// BadRequestError indicates that the client request cannot be fulfilled
type BadRequestError struct {
	Err     error
	Request *http.Request
}

// Error implements the error interface
func (e BadRequestError) Error() string {
	return fmt.Sprintf("the request '%s' cannot be fulfilled because: '%v'", e.Request.RequestURI, e.Err)
}

// ServerError is used when an unexpected situation occurred
type ServerError struct {
	Err     error
	Request *http.Request
}

// Error implements the error interface
func (e ServerError) Error() string {
	return fmt.Sprintf("the request '%s' resulted in an unexpected error: '%v'", e.Request.RequestURI, e.Err)
}

// --------------------------------------------------------------------------
// Shortcuts for commen error responses
// --------------------------------------------------------------------------

// ErrBadRequest returns a http.StatusBadRequest
func ErrBadRequest(err BadRequestError) render.Renderer {
	return &ProblemDetail{
		Type:   "about:blank",
		Title:  "the request cannot be fulfilled",
		Status: http.StatusNotFound,
		Detail: err.Error(),
	}
}

// ErrNotFound returns a http.StatusNotFound
func ErrNotFound(err NotFoundError) render.Renderer {
	return &ProblemDetail{
		Type:   "about:blank",
		Title:  "object cannot be found",
		Status: http.StatusNotFound,
		Detail: err.Error(),
	}
}

// ErrServerError returns a http.StatusInternalServerError
func ErrServerError(err ServerError) render.Renderer {
	return &ProblemDetail{
		Type:   "about:blank",
		Title:  "cannot service the request",
		Status: http.StatusNotFound,
		Detail: err.Error(),
	}
}

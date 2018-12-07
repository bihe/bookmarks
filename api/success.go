package api

import (
	"net/http"

	"github.com/go-chi/render"
)

// SuccessResponse defines a generic success status
type SuccessResponse struct {
	HTTPStatusCode int    `json:"status"`            // http response status code
	Message        string `json:"message,omitempty"` // application-level message
}

// SuccessResult created a success result
func SuccessResult(code int, message string) render.Renderer {
	return &SuccessResponse{
		HTTPStatusCode: code,
		Message:        message,
	}
}

// Render is the overloaded method for the ErrResponse
func (s *SuccessResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, s.HTTPStatusCode)
	return nil
}

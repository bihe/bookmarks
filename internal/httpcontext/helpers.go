package httpcontext

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/bihe/bookmarks-go/internal/httpcontext/header"
)

// JSON is a simple map to produce JSON serialized data
type JSON map[string]interface{}

// WriteJSON sets the HTTP status code and marshals the data to JSON
func WriteJSON(w http.ResponseWriter, code int, data map[string]interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	b, err := json.Marshal(data)
	if err != nil {
		log.SetPrefix("context.WriteJSON")
		log.Printf("could not marshal json %v\n", err)
	}
	_, err = w.Write(b)
	if err != nil {
		log.SetPrefix("context.WriteJSON")
		log.Printf("could not write bytes using http.ResponseWriter: %v\n", err)
	}
}

// Write sets the HTTP status-code and content-type for the response and writes data
func Write(w http.ResponseWriter, code int, contentType string, data []byte) {
	w.Header().Set("Content-Type", fmt.Sprintf("%s; charset=utf-8", contentType))
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	_, err := w.Write(data)
	if err != nil {
		log.SetPrefix("context.WriteJSON")
		log.Printf("could not write bytes using http.ResponseWriter: %v\n", err)
	}
}

// NegotiateError negotiates content-types, sets the status code and returns error information
func NegotiateError(w http.ResponseWriter, r *http.Request, code int, message, redirectURL string) {
	ctype := acceptContentType(r)
	switch ctype {
	case "text/html":
		if redirectURL != "" {
			http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
			return
		}
		fallthrough
	case "text/plain":
		Write(w, code, "text/plain", []byte(message))
	default:
		// default approach is to write JSON
		WriteJSON(w, code, JSON{
			"status":  code,
			"message": message,
		})
	}
}

var contentTypes = []string{"text/html", "text/plain", "application/json", "application/octet-stream"}
var defaultType = "application/json"

func acceptContentType(r *http.Request) string {
	return header.NegotiateContentType(r, contentTypes, defaultType)
}

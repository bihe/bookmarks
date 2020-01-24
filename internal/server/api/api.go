package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bihe/bookmarks/internal"
	constants "github.com/bihe/commons-go/config"
	"github.com/bihe/commons-go/cookies"
	"github.com/bihe/commons-go/errors"
	"github.com/bihe/commons-go/security"
	log "github.com/sirupsen/logrus"
)

// --------------------------------------------------------------------------
// API Interface
// --------------------------------------------------------------------------

// Bookmarks defines the available methods of the Bookmarks API
type Bookmarks interface {
	// appinfo
	HandleAppInfo(user security.User, w http.ResponseWriter, r *http.Request) error

	// wrapper methods
	Secure(f func(user security.User, w http.ResponseWriter, r *http.Request) error) http.HandlerFunc
	Call(f func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc
}

// NewBookmarksAPI creates a new instance of the Bookmarks interface
func NewBookmarksAPI(cs cookies.Settings, version internal.VersionInfo) Bookmarks {
	api := bookmarksAPI{
		VersionInfo: version,
		errRep:      errors.NewReporter(cs),
	}
	return &api
}

var _ Bookmarks = (*bookmarksAPI)(nil)

// bookmarksAPI provides the implementation of the Bookmarks Interface
type bookmarksAPI struct {
	internal.VersionInfo
	errRep *errors.ErrorReporter
}

// Secure wraps handlers to have a common signature
// a User is retrieved from the context and a possible error from the handler function is processed
func (b *bookmarksAPI) Secure(f func(user security.User, w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u := r.Context().Value(constants.User)
		if u == nil {
			log.WithField("func", "server.secure").Errorf("user is not available in context!")
			b.errRep.Negotiate(w, r, fmt.Errorf("user is not available in context"))
			return
		}
		user := r.Context().Value(constants.User).(*security.User)
		if err := f(*user, w, r); err != nil {
			log.WithField("func", "server.secure").Errorf("error during API call %v\n", err)
			b.errRep.Negotiate(w, r, err)
			return
		}
	})
}

// Call wraps handlers to have a common signature
func (b *bookmarksAPI) Call(f func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			log.WithField("func", "server.call").Errorf("error during API call %v\n", err)
			b.errRep.Negotiate(w, r, err)
			return
		}
	})
}

// --------------------------------------------------------------------------
// internal API helpers
// --------------------------------------------------------------------------

// respond converts data into appropriate responses for the client
// this can be JSON, Plaintext, ...
func (b *bookmarksAPI) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/problem+json; charset=utf-8")
	w.WriteHeader(code)
	if data != nil {
		err := json.NewEncoder(w).Encode(data)
		if err != nil {
			log.WithField("func", "server.respond").Errorf("could not marshal json %v\n", err)
			b.errRep.Negotiate(w, r, errors.ServerError{
				Err:     fmt.Errorf("could not marshal json %v", err),
				Request: r,
			})
			return
		}
	}
}

// decode parses supplied JSON payload
func (b *bookmarksAPI) decode(w http.ResponseWriter, r *http.Request, v interface{}) error {
	if r.Body == nil {
		return fmt.Errorf("no body payload available")
	}
	return json.NewDecoder(r.Body).Decode(v)
}

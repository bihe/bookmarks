package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bihe/bookmarks/internal"
	"github.com/bihe/commons-go/cookies"
	"github.com/bihe/commons-go/errors"
	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"

	constants "github.com/bihe/commons-go/config"
	"github.com/bihe/commons-go/security"
)

var cookieSettings = cookies.Settings{
	Path:   "/",
	Domain: "localhost",
	Secure: false,
	Prefix: "test",
}

var version = internal.VersionInfo{
	Version: "1",
	Build:   "2",
}

func TestGetAppInfo(t *testing.T) {
	r := chi.NewRouter()
	api := NewBookmarksAPI(cookieSettings, version)

	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), constants.User, &security.User{
				Username:    "username",
				Email:       "a.b@c.de",
				DisplayName: "displayname",
				Roles:       []string{"role"},
				UserID:      "12345",
			})
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	})
	r.Get("/appinfo", api.Secure(api.HandleAppInfo))

	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/appinfo", nil)

	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var m Meta
	err := json.Unmarshal(rec.Body.Bytes(), &m)
	if err != nil {
		t.Errorf("could not get valid json: %v", err)
	}

	assert.Equal(t, "1-2", m.Version)
}

func TestGetAppInfoNilUser(t *testing.T) {
	r := chi.NewRouter()
	api := NewBookmarksAPI(cookieSettings, version)
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), constants.User, nil)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	})
	r.Get("/appinfo", api.Secure(api.HandleAppInfo))

	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/appinfo", nil)

	r.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var p errors.ProblemDetail
	err := json.Unmarshal(rec.Body.Bytes(), &p)
	if err != nil {
		t.Errorf("could not unmarshall: %v", err)
	}
	assert.Equal(t, http.StatusInternalServerError, p.Status)
}

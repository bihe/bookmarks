package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/bihe/bookmarks/internal/store"
	"github.com/bihe/commons-go/cookies"
	"github.com/bihe/commons-go/errors"
	"github.com/bihe/commons-go/handler"
	"github.com/bihe/commons-go/security"
	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
)

// --------------------------------------------------------------------------
// test helpers
// --------------------------------------------------------------------------

const userName = "username"

var raisedError = fmt.Errorf("error")

// common components necessary for handlers
var baseHandler = handler.Handler{
	ErrRep: errors.NewReporter(cookies.Settings{
		Path:   "/",
		Domain: "localhost",
		Secure: false,
		Prefix: "test",
	}, "error"),
}

func jwtUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), security.UserKey, &security.User{
			Username:    userName,
			Email:       "a.b@c.de",
			DisplayName: "displayname",
			Roles:       []string{"role"},
			UserID:      "12345",
		})
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// --------------------------------------------------------------------------
// test-code
// --------------------------------------------------------------------------

type TestGetBookmarkByIdRepo struct {
	mockRepository
	fail bool
}

func (r *TestGetBookmarkByIdRepo) GetBookmarkById(id, username string) (store.Bookmark, error) {
	if r.fail {
		return store.Bookmark{}, raisedError
	}
	return store.Bookmark{
		DisplayName: "displayname",
		AccessCount: 1,
		ChildCount:  0,
		Created:     time.Now().UTC(),
		ID:          "ID",
		Path:        "/",
		Type:        store.Node,
		URL:         "http://url",
		UserName:    username,
	}, nil
}

func TestGetBookmarkById(t *testing.T) {
	// arrange
	r := chi.NewRouter()
	r.Use(jwtUser)
	bookmarkAPI := &BookmarksAPI{
		Handler:    baseHandler,
		Repository: &TestGetBookmarkByIdRepo{},
	}
	url := "/api/v1/bookmarks/ID"
	rec := httptest.NewRecorder()
	var bm Bookmark

	// act
	r.Get(url, bookmarkAPI.Secure(bookmarkAPI.GetBookmarkByID))
	req, _ := http.NewRequest("GET", url, nil)
	r.ServeHTTP(rec, req)

	// assert
	assert.Equal(t, http.StatusOK, rec.Code)
	if err := json.Unmarshal(rec.Body.Bytes(), &bm); err != nil {
		t.Errorf("could not unmarshal: %v", err)
	}

	assert.Equal(t, bm.ID, "ID")
	assert.Equal(t, bm.DisplayName, "displayname")
	assert.Equal(t, bm.Path, "/")
	assert.Equal(t, bm.URL, "http://url")
	assert.Equal(t, bm.Type, Node)

	// fail -------------------------------------------------------------
	bookmarkAPI.Repository = &TestGetBookmarkByIdRepo{fail: true}
	rec = httptest.NewRecorder()

	r.Get(url, bookmarkAPI.Secure(bookmarkAPI.GetBookmarkByID))
	req, _ = http.NewRequest("GET", url, nil)
	r.ServeHTTP(rec, req)

	// assert
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

type TestGetBookmarkByPathRepo struct {
	mockRepository
	fail bool
}

func (r *TestGetBookmarkByPathRepo) GetBookmarksByPath(path, username string) ([]store.Bookmark, error) {
	if r.fail {
		return make([]store.Bookmark, 0), raisedError
	}

	bm := store.Bookmark{
		DisplayName: "displayname",
		AccessCount: 1,
		ChildCount:  0,
		Created:     time.Now().UTC(),
		ID:          "ID",
		Path:        "/",
		Type:        store.Node,
		URL:         "http://url",
		UserName:    username,
	}
	var bms []store.Bookmark
	return append(bms, bm), nil
}

func TestGetBookmarkByPath(t *testing.T) {
	// arrange
	r := chi.NewRouter()
	r.Use(jwtUser)
	bookmarkAPI := &BookmarksAPI{
		Handler:    baseHandler,
		Repository: &TestGetBookmarkByPathRepo{},
	}
	reqUrl := "/api/v1/bookmarks/bypath"
	rec := httptest.NewRecorder()
	var bl BookmarkList

	q := make(url.Values)
	q.Set("path", "/")

	// act
	r.Get(reqUrl, bookmarkAPI.Secure(bookmarkAPI.GetBookmarksByPath))
	req, _ := http.NewRequest("GET", reqUrl+"?"+q.Encode(), nil)
	r.ServeHTTP(rec, req)

	// assert
	assert.Equal(t, http.StatusOK, rec.Code)
	if err := json.Unmarshal(rec.Body.Bytes(), &bl); err != nil {
		t.Errorf("could not unmarshal: %v", err)
	}

	assert.Equal(t, 1, bl.Count)
	assert.Equal(t, true, bl.Success)

	// fail -------------------------------------------------------------
	bookmarkAPI.Repository = &TestGetBookmarkByPathRepo{fail: true}
	rec = httptest.NewRecorder()

	r.Get(reqUrl, bookmarkAPI.Secure(bookmarkAPI.GetBookmarksByPath))
	req, _ = http.NewRequest("GET", reqUrl+"?"+q.Encode(), nil)
	r.ServeHTTP(rec, req)

	// assert
	assert.Equal(t, http.StatusOK, rec.Code)
	if err := json.Unmarshal(rec.Body.Bytes(), &bl); err != nil {
		t.Errorf("could not unmarshal: %v", err)
	}

	assert.Equal(t, 0, bl.Count)
	assert.Equal(t, true, bl.Success)
}

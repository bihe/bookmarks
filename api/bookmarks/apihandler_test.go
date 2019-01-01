package bookmarks_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"path"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/bihe/bookmarks/api/bookmarks"
	"github.com/bihe/bookmarks/core"
	"github.com/bihe/bookmarks/security"
	"github.com/bihe/bookmarks/store"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi"
)

func getTestConfig() core.Configuration {
	return core.Configuration{
		DB: core.Database{Dialect: "sqlite3", Connection: ":memory:"},
		Sec: core.Security{
			JwtIssuer:     "i",
			JwtSecret:     "s",
			CookieName:    "c",
			LoginRedirect: "http://locahost/redirect",
			Claim:         core.Claim{Name: "bookmarks", URL: "http://localhost", Roles: []string{"User"}},
		},
		FS: core.FileServer{
			Path:    "/tmp",
			URLPath: "/static",
		},
	}
}

type CustomClaims struct {
	Type        string   `json:"Type"`
	UserName    string   `json:"UserName"`
	Email       string   `json:"Email"`
	UserID      string   `json:"UserId"`
	DisplayName string   `json:"DisplayName"`
	Surname     string   `json:"Surname"`
	GivenName   string   `json:"GivenName"`
	Claims      []string `json:"Claims"`
	jwt.StandardClaims
}

func createToken() string {
	// Create the Claims
	claims := CustomClaims{
		"login.User",
		"a",
		"a.b@c.de",
		"1",
		"A B",
		"B",
		"A",
		[]string{"bookmarks|http://localhost|User"},
		jwt.StandardClaims{
			ExpiresAt: time.Date(2099, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
			Issuer:    "i",
		},
	}

	// Create a new token object, specifying signing method and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and get the complete encoded token as a string using the secret
	tokenString := ""
	var err error
	if tokenString, err = token.SignedString([]byte("s")); err != nil {
		return ""
	}
	return tokenString
}

func getSqliteDDL() string {
	dir, err := filepath.Abs("../../")
	if err != nil {
		return ""
	}
	return path.Join(dir, "_db/", "ddl.sql")
}

func setupApiInitDB(config core.Configuration, ddlFilePath string) *chi.Mux {
	r := chi.NewRouter()

	// configure JWT authentication and use JWT middleware
	r.Use(security.NewMiddleware(config).JWTContext)

	r.Route("/api/v1", func(r chi.Router) {
		store := store.New(config.DB.Dialect, config.DB.Connection)
		if ddlFilePath != "" {
			store.InitSchema(ddlFilePath)
		}
		r.Mount("/bookmarks", bookmarks.MountRoutes(store))
	})
	return r
}

func TestAPICreateBookmark(t *testing.T) {
	ddlFilePath := getSqliteDDL()
	if ddlFilePath == "" {
		t.Fatalf("Could not get ddl file for sqlite in memory db!")
	}

	router := setupApiInitDB(getTestConfig(), ddlFilePath)
	jwt := createToken()
	tt := []struct {
		name     string
		payload  string
		status   int
		jwt      string
		response string
	}{
		{
			name: "Create a new Bookmark",
			payload: `{
				"path":"/",
				"displayName":"Test",
				"url": "http://a.b.c.de",
				"sortOrder": 1,
				"type": "node"
			}`,
			jwt:      jwt,
			status:   http.StatusCreated,
			response: `{"status":201,"message":"bookmark item created: p:/, n:Test"`,
		},
		{
			name: "Missing folder",
			payload: `{
				"path":"/A",
				"displayName":"Test",
				"url": "http://a.b.c.de",
				"sortOrder": 1,
				"type": "node"
			}`,
			jwt:      jwt,
			status:   http.StatusBadRequest,
			response: `{"type":"about:blank","title":"the request cannot be fulfilled","status":400,"detail":"the request '' cannot be fulfilled because: cannot create item because of missing folder structure: the folder with path '/' and name 'A' does not exist"}`,
		},
		{
			name: "Invalid characters",
			payload: `{
				"path":"/",
				"displayName":"Test/",
				"url": "http://a.b.c.de",
				"sortOrder": 1,
				"type": "node"
			}`,
			jwt:      jwt,
			status:   http.StatusBadRequest,
			response: `{"type":"about:blank","title":"the request cannot be fulfilled","status":400,"detail":"the request '' cannot be fulfilled because: invalid chars in 'DisplayName'"}`,
		},
		{
			name: "Invalid path",
			payload: `{
				"path":"/a/",
				"displayName":"Test",
				"url": "http://a.b.c.de",
				"sortOrder": 1,
				"type": "node"
			}`,
			jwt:      jwt,
			status:   http.StatusBadRequest,
			response: `{"type":"about:blank","title":"the request cannot be fulfilled","status":400,"detail":"the request '' cannot be fulfilled because: a path cannot end with '/"}`,
		},
		{
			name:     "Wrong payload",
			payload:  "",
			jwt:      jwt,
			status:   http.StatusBadRequest,
			response: `{"type":"about:blank","title":"the request cannot be fulfilled","status":400,"detail":"the request '' cannot be fulfilled because: EOF"}`,
		},
		{
			name:     "No jwt auth token",
			payload:  "",
			jwt:      "",
			status:   http.StatusUnauthorized,
			response: `{"detail":"Invalid authentication, no JWT token present!","status":401,"title":"security error","type":"about:blank"}`,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/api/v1/bookmarks", strings.NewReader(tc.payload))
			req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tc.jwt))
			req.Header.Add("Content-Type", "application/json")
			router.ServeHTTP(w, req)

			if w.Code != tc.status {
				t.Fatalf("the status should be '%d' but got '%d'", tc.status, w.Code)
			}
			if strings.Index(strings.TrimSpace(w.Body.String()), tc.response) == -1 {
				t.Fatalf("expected response '%s' but got '%s'", tc.response, strings.TrimSpace(w.Body.String()))
			}
		})

	}
}

func TestAPIGetBookmarks(t *testing.T) {
	// first: create a fresh bookmark
	payload := `{
		"path":"/",
		"displayName":"Test",
		"url": "http://a.b.c.de",
		"sortOrder": 1,
		"type": "node"
	}`
	jwt := createToken()
	ddlFilePath := getSqliteDDL()
	if ddlFilePath == "" {
		t.Fatalf("Could not get ddl file for sqlite in memory db!")
	}
	router := setupApiInitDB(getTestConfig(), ddlFilePath)

	w := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "/api/v1/bookmarks", strings.NewReader(payload))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", jwt))
	req.Header.Add("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("the status should be '%d' but got '%d'", http.StatusCreated, w.Code)
	}
	var bc bookmarks.BookmarkCreatedResponse
	err = json.Unmarshal(w.Body.Bytes(), &bc)
	if err != nil {
		t.Fatal(err)
	}
	if bc.NodeID == "" {
		t.Fatalf("did not get ID from create-bookmarks-call!")
	}

	// query all bookmarks
	w = httptest.NewRecorder()
	req, err = http.NewRequest("GET", "/api/v1/bookmarks", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", jwt))
	req.Header.Add("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("the status should be '%d' but got '%d'", http.StatusOK, w.Code)
	}
	var bl bookmarks.BookmarkListResponse
	err = json.Unmarshal(w.Body.Bytes(), &bl)
	if err != nil {
		t.Fatal(err)
	}
	if bl.Count != 1 {
		t.Fatalf("expected '1' bookmarks but got '%d'", bl.Count)
	}

	// query a specific bookmark - using the NodeID from above
	w = httptest.NewRecorder()
	req, err = http.NewRequest("GET", "/api/v1/bookmarks/"+bc.NodeID, nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", jwt))
	req.Header.Add("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("the status should be '%d' but got '%d'", http.StatusOK, w.Code)
	}
	var br bookmarks.BookmarkResponse
	err = json.Unmarshal(w.Body.Bytes(), &br)
	if err != nil {
		t.Fatal(err)
	}
	if br.NodeID != bc.NodeID {
		t.Fatalf("expected ID '%s' but got ID '%s'", bc.NodeID, br.NodeID)
	}
}

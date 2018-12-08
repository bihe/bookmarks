package bookmarks_test

import (
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
	jwt "github.com/dgrijalva/jwt-go"
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

func TestAPICreateBookmark(t *testing.T) {
	ddlFilePath := getSqliteDDL()
	if ddlFilePath == "" {
		t.Fatalf("Could not get ddl file for sqlite in memory db!")
	}

	router := bookmarks.SetupAPIInitDB(getTestConfig(), ddlFilePath)
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
			response: `{"status":201,"message":"bookmark item created: p:/, n:Test"}`,
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
			response: `{"status":400,"message":"invalid request: cannot create item because of missing folder structure: the folder with path '/' and name 'A' does not exist"}`,
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
			response: `{"status":400,"message":"invalid request: invalid chars in 'DisplayName'"}`,
		},
		{
			name:     "Wrong payload",
			payload:  "",
			jwt:      jwt,
			status:   http.StatusBadRequest,
			response: `{"status":400,"message":"invalid request: EOF"}`,
		},
		{
			name:     "No jwt auth token",
			payload:  "",
			jwt:      "",
			status:   http.StatusUnauthorized,
			response: `{"message":"Invalid authentication, no JWT token present!","status":401}`,
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
				t.Errorf("the status should be '%d' but got '%d'", tc.status, w.Code)
			}
			if strings.TrimSpace(w.Body.String()) != tc.response {
				t.Errorf("expected response '%s' but got '%s'", tc.response, strings.TrimSpace(w.Body.String()))
			}
		})

	}
}

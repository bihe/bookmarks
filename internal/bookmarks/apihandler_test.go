package bookmarks

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestCreateBookmark(t *testing.T) {
	router := SetupRouter(os.Getenv("APPLICATION_BASE_PATH"), os.Getenv("APPLICATION_CONFIG_FILE"))
	jwt := os.Getenv("JWT_TOKEN")

	tt := []struct {
		name     string
		payload  string
		status   int
		response string
	}{
		{
			name: "Create a new Bookmark",
			payload: `{
				"path":"/A/B/C",
				"displayName":"Test",
				"url": "http://a.b.c.de",
				"sortOrder": 1,
				"itemType": "node"
			}`,
			status:   http.StatusCreated,
			response: `{"status":201,"message":"bookmark item created: /A/B/C/Test"}`,
		},
		{
			name:     "Wrong payload",
			payload:  "",
			status:   http.StatusBadRequest,
			response: `{"status":400,"message":"invalid request: EOF"}`,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/api/v1/bookmarks", strings.NewReader(tc.payload))
			req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", jwt))
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

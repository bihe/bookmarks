package bookmarks

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateBookmark(t *testing.T) {
	router := SetupRouter(os.Getenv("APPLICATION_BASE_PATH"), os.Getenv("APPLICATION_CONFIG_FILE"), false)
	jwt := os.Getenv("JWT_TOKEN")

	w := httptest.NewRecorder()

	var payload = `{
		"path":"/A/B/C",
		"displayName":"Test",
		"url": "http://a.b.c.de",
		"sortOrder": 1,
		"itemType": "node"
	}`

	req, _ := http.NewRequest("POST", "/api/v1/bookmarks", strings.NewReader(payload))
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", jwt))
	req.Header.Add("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, 201, w.Code)
	assert.Equal(t, `{"message":"bookmark item created: /A/B/C/Test","status":201}`, w.Body.String())
}

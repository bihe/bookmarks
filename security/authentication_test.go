package security_test

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bihe/bookmarks/security"
)

// GetTestHandler returns a http.HandlerFunc for testing http middleware
func GetTestHandler() http.HandlerFunc {
	fn := func(rw http.ResponseWriter, req *http.Request) {
		panic("test entered test handler, this should not happen")
	}
	return http.HandlerFunc(fn)
}

func TestAuthenticationHandler(t *testing.T) {
	o := security.AuthOptions{
		JwtIssuer:     "i",
		JwtSecret:     "s",
		CookieName:    "c",
		RedirectURL:   "http://locahost/redirect",
		RequiredClaim: security.Claim{Name: "bookmarks", URL: "http://localhost", Roles: []string{"User"}},
	}
	j := security.JwtMiddleware{Options: o}
	ts := httptest.NewServer(j.JWTContext(GetTestHandler()))
	defer ts.Close()

	var u bytes.Buffer
	u.WriteString(string(ts.URL))
	u.WriteString("/")

	client := &http.Client{}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		t.Errorf("cannot created a client request: %v", err)
		return
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"))
	res, err := client.Do(req)
	if err != nil {
		t.Errorf("error callind middleware: %v", err)
		return
	}
	if res.StatusCode != http.StatusUnauthorized {
		t.Errorf("invalid status-code returned; wanted %d got %d", http.StatusUnauthorized, res.StatusCode)
	}
}

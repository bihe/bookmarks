package security

import (
	"testing"
)

func TestAuthorization(t *testing.T) {
	reqClaim := Claim{Name: "test", URL: "http://a.b.c.de/f", Roles: []string{"user", "admin"}}
	var claims []string

	claims = append(claims, "a|http://1.com|nothing")
	claims = append(claims, "test|http://a.b.c.de/f|user")

	if _, err := Authorize(reqClaim, claims); err != nil {
		t.Error("Authorization failed.", err)
	}

	claims = claims[0:1]
	if _, err := Authorize(reqClaim, claims); err == nil {
		t.Error("Authorization should fail.", err)
	}

	claims = nil
	claims = append(claims, "a|http://1.com|nothing")
	claims = append(claims, "test|http://a.b.c.de/f|admin")

	if _, err := Authorize(reqClaim, claims); err != nil {
		t.Error("Authorization failed.", err)
	}
}

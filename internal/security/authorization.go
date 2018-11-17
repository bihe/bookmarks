package security

import (
	"fmt"
	"net/url"
	"strings"
)

// Claim defines the authorization requiremenets
type Claim struct {
	// Name of the applicatiion
	Name string
	// URL of the application
	URL string
	// Role possible roles
	Roles []string
}

// Authorize validates the given claims and verifies if
// they match the required claim
// a claim entriy is in the form "name|url|role"
func Authorize(required Claim, claims []string) (roles []string, err error) {
	for _, claim := range claims {
		c := split(claim)
		ok, _ := compareURL(required.URL, c.URL)
		if required.Name == c.Name && matchRole(c.Roles, required.Roles) && ok {
			return c.Roles, nil
		}
	}
	return roles, fmt.Errorf("supplied claims are not sufficient")
}

func matchRole(a []string, b []string) bool {
	for _, r := range a {
		for _, s := range b {
			if s == r {
				return true
			}
		}
	}
	return false
}

func split(claim string) *Claim {
	parts := strings.Split(claim, "|")
	if len(parts) == 3 {
		return &Claim{Name: parts[0], URL: parts[1], Roles: []string{parts[2]}}
	}
	return &Claim{}
}

func compareURL(a, b string) (bool, error) {
	var (
		urlA *url.URL
		urlB *url.URL
		err  error
	)
	if urlA, err = url.Parse(a); err != nil {
		return false, err
	}
	if urlB, err = url.Parse(b); err != nil {
		return false, err
	}
	if urlA.Scheme != urlB.Scheme || urlA.Port() != urlB.Port() || urlA.Host != urlB.Host {
		return false, fmt.Errorf("The urls do not match: '%s vs. %s'", urlA, urlB)
	}
	if urlA.Path != urlB.Path {
		return false, fmt.Errorf("The path of the urls does not match: '%s vs. %s'", urlA.Path, urlB.Path)
	}
	return true, nil
}

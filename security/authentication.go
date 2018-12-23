package security

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/bihe/bookmarks/core"
	"github.com/bihe/bookmarks/security/httpcontext"
)

// AuthOptions defines presets for the Authentication handler
// by the default the JWT token is fetched from the Authentication header
// as a fallback it is possible to fetch the token from a specific cookie
type AuthOptions struct {
	// JwtSecret is the jwt signing key
	JwtSecret string
	// JwtIssuer specifies identifies the principal that issued the token
	JwtIssuer string
	// CookieName spedifies the HTTP cookie holding the token
	CookieName string
	// RequiredClaim to access the application
	RequiredClaim Claim
	// RedirectURL forwards the request to an external authentication service
	RedirectURL string
}

// User is the authenticated principal extracted from the JWT token
type User struct {
	Username    string
	Roles       []string
	Email       string
	UserID      string
	DisplayName string
}

// JwtMiddleware is responsible for JWT authentication and authorization
type JwtMiddleware struct {
	Options AuthOptions
}

// NewMiddleware created a new instance using the supplied config options
func NewMiddleware(config core.Configuration) *JwtMiddleware {
	return &JwtMiddleware{
		Options: AuthOptions{
			CookieName: config.Sec.CookieName,
			JwtIssuer:  config.Sec.JwtIssuer,
			JwtSecret:  config.Sec.JwtSecret,
			RequiredClaim: Claim{
				Name:  config.Sec.Claim.Name,
				URL:   config.Sec.Claim.URL,
				Roles: config.Sec.Claim.Roles,
			},
			RedirectURL: config.Sec.LoginRedirect,
		},
	}
}

// JWTContext parses provided information from the request and populates user-data
// in the request or denies access if required data is missing
func (jwt *JwtMiddleware) JWTContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var token string
		var err error
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" {
			token = strings.Replace(authHeader, "Bearer ", "", 1)
		}
		if token == "" {
			// fallback to get the token via the cookie
			var cookie *http.Cookie
			if cookie, err = r.Cookie(jwt.Options.CookieName); err != nil {
				// neither the header nor the cookie supplied a jwt token
				httpcontext.NegotiateError(w, r, http.StatusUnauthorized, "Invalid authentication, no JWT token present!", jwt.Options.RedirectURL)
				return
			}
			token = cookie.Value
		}
		var payload JwtTokenPayload
		if payload, err = ParseJwtToken(token, jwt.Options.JwtSecret, jwt.Options.JwtIssuer); err != nil {
			log.Printf("Could not decode the JWT token payload: %s", err)
			httpcontext.NegotiateError(w, r, http.StatusUnauthorized, fmt.Sprintf("Invalid authentication, could not parse the JWT token: %v", err), jwt.Options.RedirectURL)
			return
		}
		var roles []string
		if roles, err = Authorize(jwt.Options.RequiredClaim, payload.Claims); err != nil {
			log.Printf("Insufficient permissions to access the resource: %s", err)
			httpcontext.NegotiateError(w, r, http.StatusForbidden, fmt.Sprintf("Invalid authorization: %v", err), jwt.Options.RedirectURL)
			return
		}
		user := &User{
			DisplayName: payload.DisplayName,
			Email:       payload.Email,
			Roles:       roles,
			UserID:      payload.UserID,
			Username:    payload.UserName,
		}

		ctx := context.WithValue(r.Context(), core.ContextUser, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

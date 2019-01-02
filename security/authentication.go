package security

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/bihe/bookmarks/cache"
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
	// CacheDuration defines the duration to cache the JWT token result
	CacheDuration string
}

// User is the authenticated principal extracted from the JWT token
type User struct {
	Username    string
	Roles       []string
	Email       string
	UserID      string
	DisplayName string
}

func (u User) serialize() ([]byte, error) {
	b, err := json.Marshal(u)
	return b, err
}

func (u *User) deserialize(b []byte) error {
	err := json.Unmarshal(b, u)
	return err
}

// JwtMiddleware is responsible for JWT authentication and authorization
type JwtMiddleware struct {
	options AuthOptions
	cache   *cache.MemoryCache
}

// NewMiddleware creates a new instance using the supplied config options
func NewMiddleware(config core.Configuration) *JwtMiddleware {
	return &JwtMiddleware{
		options: AuthOptions{
			CookieName: config.Sec.CookieName,
			JwtIssuer:  config.Sec.JwtIssuer,
			JwtSecret:  config.Sec.JwtSecret,
			RequiredClaim: Claim{
				Name:  config.Sec.Claim.Name,
				URL:   config.Sec.Claim.URL,
				Roles: config.Sec.Claim.Roles,
			},
			RedirectURL:   config.Sec.LoginRedirect,
			CacheDuration: config.Sec.CacheDuration,
		},
		cache: cache.NewCache(),
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
			if cookie, err = r.Cookie(jwt.options.CookieName); err != nil {
				// neither the header nor the cookie supplied a jwt token
				httpcontext.NegotiateError(w, r, http.StatusUnauthorized, "Invalid authentication, no JWT token present!", jwt.options.RedirectURL)
				return
			}
			token = cookie.Value
		}

		// to speed up processing use the cache for token lookups
		var user User
		u := jwt.cache.Get(token)
		if u != nil {
			err = user.deserialize(u)
			if err == nil {
				ctx := context.WithValue(r.Context(), core.ContextUser, &user)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			log.Printf("Could not use cached value of user, continue validating JWT token: %v.", err)
		}

		var payload JwtTokenPayload
		if payload, err = ParseJwtToken(token, jwt.options.JwtSecret, jwt.options.JwtIssuer); err != nil {
			log.Printf("Could not decode the JWT token payload: %s", err)
			httpcontext.NegotiateError(w, r, http.StatusUnauthorized, fmt.Sprintf("Invalid authentication, could not parse the JWT token: %v", err), jwt.options.RedirectURL)
			return
		}
		var roles []string
		if roles, err = Authorize(jwt.options.RequiredClaim, payload.Claims); err != nil {
			log.Printf("Insufficient permissions to access the resource: %s", err)
			httpcontext.NegotiateError(w, r, http.StatusForbidden, fmt.Sprintf("Invalid authorization: %v", err), jwt.options.RedirectURL)
			return
		}

		user = User{
			DisplayName: payload.DisplayName,
			Email:       payload.Email,
			Roles:       roles,
			UserID:      payload.UserID,
			Username:    payload.UserName,
		}

		u, err = user.serialize()
		if err != nil {
			log.Printf("Could not marshall the User object for caching: %v", err)
		} else {
			d, err := time.ParseDuration(jwt.options.CacheDuration)
			if err == nil {
				jwt.cache.Set(token, u, d)
			} else {
				log.Printf("Could not cache User object because of duration parsing error: %v", err)
			}
		}

		ctx := context.WithValue(r.Context(), core.ContextUser, &user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

package security

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/bihe/bookmarks-go/internal/conf"
	"github.com/bihe/bookmarks-go/internal/context"

	"github.com/gin-gonic/gin"
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
	Role        []string
	Email       string
	UserID      string
	DisplayName string
}

// JwtAuth parses provided information from the request and populates user-data
// in the request or denies access if required data is missing
func JwtAuth(options AuthOptions) gin.HandlerFunc {
	return func(c *gin.Context) {
		var token string
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			token = strings.Replace(authHeader, "Bearer ", "", 1)
		}
		if token == "" {
			// fallback to get the token via the cookie
			var err error
			if token, err = c.Cookie(options.CookieName); err != nil {
				log.Printf("Could not get the JWT token from the cookie: %s", err)
			}
		}
		if token == "" {
			// neither the header nor the cookie supplied a jwt token
			context.AbortAndRedirect(c, http.StatusUnauthorized, "Invalid authentication. No JWT token present!", options.RedirectURL)
			return
		}
		var payload JwtTokenPayload
		var err error
		if payload, err = ParseJwtToken(token, options.JwtSecret, options.JwtIssuer); err != nil {
			log.Printf("Could not decode the JWT token payload: %s", err)
			context.AbortAndRedirect(c, http.StatusUnauthorized, fmt.Sprintf("Invalid authentication. Could not parse the JWT token! %s", err), options.RedirectURL)
			return
		}
		var roles []string
		if roles, err = Authorize(options.RequiredClaim, payload.Claims); err != nil {
			log.Printf("Insufficient permissions to access the resource: %s", err)
			context.AbortAndRedirect(c, http.StatusForbidden, fmt.Sprintf("Invalid authorization. %s", err), options.RedirectURL)
			return
		}
		c.Set(conf.ContextUser, &User{
			DisplayName: payload.DisplayName,
			Email:       payload.Email,
			Role:        roles,
			UserID:      payload.UserID,
			Username:    payload.UserName,
		})

		c.Next()
	}
}

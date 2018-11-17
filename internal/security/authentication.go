package security

import (
	"fmt"
	"log"
	"net/http"
	"strings"

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
			abort(c, http.StatusUnauthorized, "Invalid authentication. No JWT token present!", options)
			return
		}
		var payload JwtTokenPayload
		var err error
		if payload, err = ParseJwtToken(token, options.JwtSecret, options.JwtIssuer); err != nil {
			log.Printf("Could not decode the JWT token payload: %s", err)
			abort(c, http.StatusUnauthorized, fmt.Sprintf("Invalid authentication. Could not parse the JWT token! %s", err), options)
			return
		}
		var roles []string
		if roles, err = Authorize(options.RequiredClaim, payload.Claims); err != nil {
			log.Printf("Insufficient permissions to access the resource: %s", err)
			abort(c, http.StatusForbidden, fmt.Sprintf("Invalid authorization. %s", err), options)
			return
		}
		c.Set("User", User{
			DisplayName: payload.DisplayName,
			Email:       payload.Email,
			Role:        roles,
			UserID:      payload.UserID,
			Username:    payload.UserName,
		})

		c.Next()
	}
}

func abort(c *gin.Context, status int, message string, options AuthOptions) {
	switch c.NegotiateFormat(gin.MIMEHTML, gin.MIMEJSON, gin.MIMEPlain) {
	case gin.MIMEJSON:
		c.AbortWithStatusJSON(status, gin.H{
			"status":  status,
			"message": message,
		})
	case gin.MIMEHTML:
		c.Redirect(http.StatusTemporaryRedirect, options.RedirectURL)
		c.Abort()
	case gin.MIMEPlain:
		c.String(status, message)
		c.Abort()
	default:
		c.AbortWithStatusJSON(status, gin.H{
			"status":  status,
			"message": message,
		})
	}
}

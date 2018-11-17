package security

import (
	"fmt"

	"github.com/dgrijalva/jwt-go"
)

// JwtTokenPayload is the parsed contents of the given token
type JwtTokenPayload struct {
	Type        string
	UserName    string
	Email       string
	Claims      []string
	UserID      string `json:"UserId"`
	DisplayName string
	Surname     string
	GivenName   string
	jwt.StandardClaims
}

// ParseJwtToken parses, validates and extracts data from a jwt token
func ParseJwtToken(token, tokenSecret, issuer string) (JwtTokenPayload, error) {
	parsed, err := jwt.ParseWithClaims(token, &JwtTokenPayload{}, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return JwtTokenPayload{}, err
	}
	if parsed.Valid {
		claims, ok := parsed.Claims.(*JwtTokenPayload)
		if ok {
			if validIssuer := claims.StandardClaims.VerifyIssuer(issuer, true); !validIssuer {
				return JwtTokenPayload{}, fmt.Errorf("Invalid issued created the token")
			}
			return *claims, nil
		}
	} else if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			return JwtTokenPayload{}, fmt.Errorf("That's not even a token")
		} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			// Token is either expired or not active yet
			return JwtTokenPayload{}, fmt.Errorf("Token is expired or not valid yet")
		}
	}
	return JwtTokenPayload{}, fmt.Errorf("could not parse JWT token")
}

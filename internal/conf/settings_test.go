package conf

import (
	"strings"
	"testing"
)

const config = `{
    "security": {
        "jwtIssuer": "login.binggl.net",
        "jwtSecret": "secret",
	"cookieName": "login_token",
	"loginRedirect": "https://login.url.com",
        "claim": {
            "name": "bookmarks",
            "url": "http://localhost:3000",
            "roles": ["User", "Admin"]
	}
    },
    "database": {
	"connectionString": "./bookmarks.db"
    }
}`

// TestConfigReader reads config settings from json
func TestConfigReader(t *testing.T) {
	reader := strings.NewReader(config)
	config, err := Defaults(reader)
	if err != nil {
		t.Error("Could not read.", err)
	}

	if config.Sec.JwtSecret != "secret" || config.Sec.Claim.Name != "bookmarks" || config.Sec.LoginRedirect != "https://login.url.com" || config.DB.Connection != "./bookmarks.db" {
		t.Error("Config values not read!")
	}
}

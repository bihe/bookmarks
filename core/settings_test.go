package core

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
	"connectionString": "./bookmarks.db",
	"dialect": "sqlite"
    },
    "logging": {
	"logPrefix": "prefix",
	"rollingFileLogger": {
		"filePath": "/temp/file",
		"maxFileSize": 100,
		"numberOfMaxBackups": 4,
		"maxAge": 7,
		"compressFile": false
	}
    }
}`

// TestConfigReader reads config settings from json
func TestConfigReader(t *testing.T) {
	reader := strings.NewReader(config)
	config, err := Settings(reader)
	if err != nil {
		t.Error("Could not read.", err)
	}

	if config.Sec.JwtSecret != "secret" || config.Sec.Claim.Name != "bookmarks" || config.Sec.LoginRedirect != "https://login.url.com" || config.DB.Connection != "./bookmarks.db" || config.DB.Dialect != "sqlite" {
		t.Error("Config values not read!")
	}

	if config.Log.Prefix != "prefix" {
		t.Error("Could not read logging settings!")
	}

	if config.Log.Rolling.FilePath != "/temp/file" {
		t.Error("Could not read rolling file settings!")
	}

	if config.Log.Rolling.MaxBackups != 4 {
		t.Error("Could not read rolling file settings!")
	}
}

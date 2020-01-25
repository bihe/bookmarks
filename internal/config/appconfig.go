// Package config defines the customization/configuration of the application
package config

import (
	"encoding/json"
	"io"
	"io/ioutil"
)

// AppConfig holds the application configuration
type AppConfig struct {
	Sec         Security           `json:"security"`
	DB          Database           `json:"database"`
	Log         LogConfig          `json:"logging"`
	Cookies     ApplicationCookies `json:"cookies"`
	ErrorPath   string             `json:"errorPath"`
	StartURL    string             `json:"startUrl"`
	Environment string             `json:"environment"`
}

// Security settings for the application
type Security struct {
	JwtIssuer     string `json:"jwtIssuer"`
	JwtSecret     string `json:"jwtSecret"`
	CookieName    string `json:"cookieName"`
	LoginRedirect string `json:"loginRedirect"`
	Claim         Claim  `json:"claim"`
	CacheDuration string `json:"cacheDuration"`
}

// Database defines the connection string
type Database struct {
	ConnStr string `json:"connectionString"`
}

// Claim defines the required claims
type Claim struct {
	Name  string   `json:"name"`
	URL   string   `json:"url"`
	Roles []string `json:"roles"`
}

// LogConfig is used to define settings for the logging process
type LogConfig struct {
	FilePath    string `json:"filePath"`
	RequestPath string `json:"requestPath"`
	LogLevel    string `json:"logLevel"`
}

// ApplicationCookies defines values for cookies
type ApplicationCookies struct {
	Domain string `json:"domain"`
	Path   string `json:"path"`
	Secure bool   `json:"secure"`
	Prefix string `json:"prefix"`
}

// GetSettings returns application configuration values
func GetSettings(r io.Reader) (*AppConfig, error) {
	var (
		c    AppConfig
		cont []byte
		err  error
	)
	if cont, err = ioutil.ReadAll(r); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(cont, &c); err != nil {
		return nil, err
	}

	return &c, nil
}

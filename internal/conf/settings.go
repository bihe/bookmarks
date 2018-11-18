package conf

import (
	"encoding/json"
	"io"
	"io/ioutil"
)

// Config holds the application configuration
type Config struct {
	Sec Security `json:"security"`
	DB  Database `json:"database"`
}

// Security settings for the application
type Security struct {
	JwtIssuer     string `json:"jwtIssuer"`
	JwtSecret     string `json:"jwtSecret"`
	CookieName    string `json:"cookieName"`
	LoginRedirect string `json:"loginRedirect"`
	Claim         Claim  `json:"claim"`
}

// Database defines the connection string
type Database struct {
	Connection string `json:"connectionString"`
}

// Claim defines the required claims
type Claim struct {
	Name  string   `json:"name"`
	URL   string   `json:"url"`
	Roles []string `json:"roles"`
}

// Defaults returns application configuration values
func Defaults(r io.Reader) (*Config, error) {
	var (
		c    Config
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

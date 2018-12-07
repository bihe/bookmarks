package core

import (
	"encoding/json"
	"io"
	"io/ioutil"
)

// Configuration holds the application configuration
type Configuration struct {
	Sec Security  `json:"security"`
	DB  Database  `json:"database"`
	Log LogConfig `json:"logging"`
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
	Dialect    string `json:"dialect"`
}

// Claim defines the required claims
type Claim struct {
	Name  string   `json:"name"`
	URL   string   `json:"url"`
	Roles []string `json:"roles"`
}

// LogConfig is used to define settings for the logging process
type LogConfig struct {
	Prefix  string        `json:"logPrefix"`
	Rolling RollingLogger `json:"rollingFileLogger"`
}

// RollingLogger defines settings to use for rolling file loggers
type RollingLogger struct {
	FilePath   string `json:"filePath"` // in megabytes
	MaxSize    int    `json:"maxFileSize"`
	MaxBackups int    `json:"numberOfMaxBackups"`
	MaxAge     int    `json:"maxAge"` // days
	Compress   bool   `json:"compressFile"`
}

// Settings returns application configuration values
func Settings(r io.Reader) (*Configuration, error) {
	var (
		c    Configuration
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

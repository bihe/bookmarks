package main

import (
	"os"
	"path"
	"path/filepath"

	"github.com/bihe/bookmarks/internal/conf"
	"github.com/bihe/bookmarks/internal/handler"
	"github.com/bihe/bookmarks/internal/request"
	"github.com/bihe/bookmarks/internal/security"
	"github.com/gin-gonic/gin"
)

func serverDefaults() (host, port string, traceLog bool) {

	port = getOrDefault("GIN_SERVER_PORT", "3000")
	host = getOrDefault("GIN_SEVER_HOST", "localhost")
	traceLog = false
	if getOrDefault("GIN_REQUEST_LOG", "0") == "1" {
		traceLog = true
	}
	return
}

func getOrDefault(env, def string) string {
	value := os.Getenv(env)
	if value == "" {
		value = def
	}
	return value
}

func main() {
	// either get the server host:port from the environment
	// or use sensible defaults
	host, port, traceLog := serverDefaults()
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic("Could not get the application basepath!")
	}
	configFile, err := os.Open(path.Join(dir, "application.json"))
	if err != nil {
		panic("Specified config file 'application.json' missing!")
	}
	defer configFile.Close()

	config, err := conf.Defaults(configFile)
	if err != nil {
		panic("No config values available to start the server. Missing config.json file!")
	}

	r := gin.Default()
	r.Use(request.Trace(traceLog))

	// group all api calls under a versioned-prefix
	api := r.Group("/api/v1")

	api.Use(security.JwtAuth(security.AuthOptions{
		CookieName: config.Sec.CookieName,
		JwtIssuer:  config.Sec.JwtIssuer,
		JwtSecret:  config.Sec.JwtSecret,
		RequiredClaim: security.Claim{
			Name:  config.Sec.Claim.Name,
			URL:   config.Sec.Claim.URL,
			Roles: config.Sec.Claim.Roles,
		},
		RedirectURL: config.Sec.LoginRedirect,
	}))
	{
		api.GET("/bookmarks", handler.GetAllBookmarks)
	}

	r.Run(host + ":" + port)
}

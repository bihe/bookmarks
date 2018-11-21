package main

import (
	"os"
	"path"
	"path/filepath"

	"github.com/bihe/bookmarks-go/internal/bookmarks"
	"github.com/bihe/bookmarks-go/internal/conf"
	"github.com/bihe/bookmarks-go/internal/request"
	"github.com/bihe/bookmarks-go/internal/security"
	"github.com/bihe/bookmarks-go/internal/store"

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

	jwtOptions := security.AuthOptions{
		CookieName: config.Sec.CookieName,
		JwtIssuer:  config.Sec.JwtIssuer,
		JwtSecret:  config.Sec.JwtSecret,
		RequiredClaim: security.Claim{
			Name:  config.Sec.Claim.Name,
			URL:   config.Sec.Claim.URL,
			Roles: config.Sec.Claim.Roles,
		},
		RedirectURL: config.Sec.LoginRedirect,
	}

	api.Use(security.JwtAuth(jwtOptions), store.InUnitOfWork(config.DB.Connection), bookmarks.CheckContext())
	{
		app := &bookmarks.Controller{}
		api.OPTIONS("/bookmarks", func(c *gin.Context) {})

		api.GET("/bookmarks/__init", app.DebugInitBookmarks)
		api.GET("/bookmarks", app.GetAllBookmarks)
	}

	r.Run(host + ":" + port)
}

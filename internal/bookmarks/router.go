package bookmarks

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/bihe/bookmarks-go/internal/conf"
	"github.com/bihe/bookmarks-go/internal/infrastructure"
	"github.com/bihe/bookmarks-go/internal/security"
	"github.com/bihe/bookmarks-go/internal/store"
	"github.com/gin-gonic/gin"
)

// SetupRouter configures the API
func SetupRouter(appBasePath, configFileName string, traceLog bool) *gin.Engine {
	dir, err := filepath.Abs(appBasePath)
	if err != nil {
		panic("Could not get the application basepath!")
	}
	configFile, err := os.Open(path.Join(dir, configFileName))
	if err != nil {
		panic(fmt.Sprintf("Specified config file '%s' missing!", configFileName))
	}
	defer configFile.Close()

	config, err := conf.Settings(configFile)
	if err != nil {
		panic("No config values available to start the server. Missing config.json file!")
	}

	r := gin.Default()
	r.Use(infrastructure.Trace(traceLog))

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

	api.Use(security.JwtAuth(jwtOptions), store.InUnitOfWork(config.DB.Connection), infrastructure.CheckContext())
	{
		bookmarks := api.Group("/bookmarks")
		{
			app := &BookmarkController{}
			bookmarks.OPTIONS("", func(c *gin.Context) {})
			bookmarks.GET("__init", app.DebugINIT)
			bookmarks.GET("/gt/:name", app.GetByID)
			bookmarks.GET("", app.GetAll)
			bookmarks.POST("", app.Create)

		}
	}

	return r
}

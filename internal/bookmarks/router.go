package bookmarks

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/bihe/bookmarks-go/internal/conf"
	"github.com/bihe/bookmarks-go/internal/security"
	"github.com/bihe/bookmarks-go/internal/store"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

// SetupRouter configures the API
func SetupRouter(appBasePath, configFileName string) *chi.Mux {
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

	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	r.Route("/api/v1", func(r chi.Router) {
		j := &security.JwtMiddleware{
			Options: jwtOptions,
		}
		u := &store.UnitOfWorkMiddleware{
			ConnStr:   config.DB.Connection,
			DbDialect: config.DB.Dialect,
		}
		r.Use(j.JWTContext)
		r.Use(u.UnitOfWorkContext)

		b := &BookmarkController{}
		r.Route("/bookmarks", func(r chi.Router) {
			r.Get("/", b.GetAll)
			r.Post("/", b.Create)
			r.Get("/{NodeID}", b.GetByID)
		})

	})

	return r
}

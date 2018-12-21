package bookmarks

import (
	"net/http"
	"strings"
	"time"

	"github.com/bihe/bookmarks/core"
	"github.com/bihe/bookmarks/security"
	"github.com/bihe/bookmarks/store"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

// SetupAPIInitDB configures the API and inits the DB - execute ddl.sql
func SetupAPIInitDB(config core.Configuration, ddlFilePath string) *chi.Mux {
	return setupAPI(config, ddlFilePath)
}

// SetupAPI configures the API
func SetupAPI(config core.Configuration) *chi.Mux {
	return setupAPI(config, "")
}

func setupAPI(config core.Configuration, ddlFilePath string) *chi.Mux {
	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.DefaultCompress)
	r.Use(middleware.RedirectSlashes)
	r.Use(middleware.Recoverer)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	// configure JWT authentication and use JWT middleware
	j := &security.JwtMiddleware{
		Options: security.AuthOptions{
			CookieName: config.Sec.CookieName,
			JwtIssuer:  config.Sec.JwtIssuer,
			JwtSecret:  config.Sec.JwtSecret,
			RequiredClaim: security.Claim{
				Name:  config.Sec.Claim.Name,
				URL:   config.Sec.Claim.URL,
				Roles: config.Sec.Claim.Roles,
			},
			RedirectURL: config.Sec.LoginRedirect,
		},
	}
	r.Use(j.JWTContext)

	// setup static file serving
	fileServer(r, config.FS.URLPath, http.Dir(config.FS.Path))

	r.Route("/api/v1", func(r chi.Router) {
		j := &security.JwtMiddleware{
			Options: security.AuthOptions{
				CookieName: config.Sec.CookieName,
				JwtIssuer:  config.Sec.JwtIssuer,
				JwtSecret:  config.Sec.JwtSecret,
				RequiredClaim: security.Claim{
					Name:  config.Sec.Claim.Name,
					URL:   config.Sec.Claim.URL,
					Roles: config.Sec.Claim.Roles,
				},
				RedirectURL: config.Sec.LoginRedirect,
			},
		}
		r.Use(j.JWTContext)
		uow := store.New(config.DB.Dialect, config.DB.Connection)
		if ddlFilePath != "" {
			uow.InitSchema(ddlFilePath)
		}
		b := &BookmarkAPI{
			uow: uow,
		}
		r.Route("/bookmarks", func(r chi.Router) {
			r.Get("/", b.GetAll)
			r.Post("/", b.Create)
			r.Put("/", b.Update)
			r.Get("/{NodeID}", b.GetByID)
			r.Get("/path", b.FindByPath)
			r.Delete("/{NodeID}", b.Delete)
			r.Delete("/{NodeID}/{Force}", b.Delete)
			r.Get("/search", b.FindByName)
		})
	})
	return r
}

func fileServer(r chi.Router, path string, root http.FileSystem) {
	if path == "" {
		panic("no path for fileServer defined!")
	}
	if strings.ContainsAny(path, "{}*") {
		panic("fileServer does not permit URL parameters.")
	}
	fs := http.StripPrefix(path, http.FileServer(root))
	// add a slash to the end of the path
	if path != "/" && path[len(path)-1] != '/' {
		path += "/"
	}
	path += "*"
	r.Get(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	}))
}

package bookmarks

import (
	"time"

	"github.com/bihe/bookmarks-go/bookmarks/conf"
	"github.com/bihe/bookmarks-go/bookmarks/security"
	"github.com/bihe/bookmarks-go/bookmarks/store"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

// SetupAPIInitDB configures the API and inits the DB - execute ddl.sql
func SetupAPIInitDB(config conf.Configuration, ddlFilePath string) *chi.Mux {
	return setupAPI(config, ddlFilePath)
}

// SetupAPI configures the API
func SetupAPI(config conf.Configuration) *chi.Mux {
	return setupAPI(config, "")
}

func setupAPI(config conf.Configuration, ddlFilePath string) *chi.Mux {
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
		uow := store.NewUnitOfWork(config.DB.Dialect, config.DB.Connection)
		if ddlFilePath != "" {
			uow.InitSchema(ddlFilePath)
		}
		b := &BookmarkController{
			uow: uow,
		}
		r.Route("/bookmarks", func(r chi.Router) {
			r.Get("/", b.GetAll)
			r.Post("/", b.Create)
			r.Put("/", b.Update)
			r.Get("/{NodeID}", b.GetByID)
			r.Get("/path", b.FindByPath)
		})
	})
	return r
}

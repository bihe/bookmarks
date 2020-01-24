package server

import (
	"net/http"
	"path/filepath"
	"time"

	"github.com/bihe/commons-go/security"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

// routes performs setup of middlewares and API handlers
func (s *Server) routes() {
	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.DefaultCompress)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	// serving static content
	serveStaticFile(r, "/favicon.ico", filepath.Join(s.basePath, "./assets/favicon.ico"))
	serveStaticDir(r, "/assets", http.Dir(filepath.Join(s.basePath, "./assets")))

	// this group "indicates" that all routes within this group use the JWT authentication
	r.Group(func(r chi.Router) {
		// authenticate and authorize users via JWT
		r.Use(security.NewJwtMiddleware(s.jwtOpts, s.cookieSettings).JwtContext)

		// group API methods together
		r.Route("/api/v1", func(r chi.Router) {
			r.Get("/appinfo", s.api.Secure(s.api.HandleAppInfo))
			// r.Get("/sites", s.api.Secure(s.api.HandleGetSites))
			// r.Post("/sites", s.api.Secure(s.api.HandleSaveSites))
			// r.Get("/sites/users/{siteName}", s.api.Secure((s.api.HandleGetUsersForSite)))
		})
		// the SPA
		serveStaticDir(r, "/ui", http.Dir(filepath.Join(s.basePath, "./assets/ui")))

		// swagger
		serveStaticDir(r, "/swagger", http.Dir(filepath.Join(s.basePath, "./assets/swagger")))
	})

	r.Get("/", http.RedirectHandler("/ui", http.StatusMovedPermanently).ServeHTTP)
	s.router = r
}

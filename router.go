package main

import (
	"net/http"
	"strings"
	"time"

	"github.com/bihe/bookmarks/api"
	"github.com/bihe/bookmarks/api/bookmarks"
	"github.com/bihe/bookmarks/core"
	"github.com/bihe/bookmarks/security"
	"github.com/bihe/bookmarks/store"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

// SetupAPI configures the API
func SetupAPI(config core.Configuration) *chi.Mux {
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
	r.Use(security.NewMiddleware(config).JWTContext)

	// setup static file serving
	fileServer(r, config.FS.URLPath, http.Dir(config.FS.Path))

	r.Route("/api/v1", func(r chi.Router) {
		store := store.New(config.DB.Dialect, config.DB.Connection)
		r.Mount("/bookmarks", bookmarks.MountRoutes(store))
		r.Mount("/appinfo", api.MountRoutes(Version, Build))
		
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

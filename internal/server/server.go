// Package server implements the API backend of the bookmark application
package server

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/bihe/commons-go/cookies"
	"github.com/bihe/commons-go/security"

	"github.com/bihe/bookmarks/internal/config"
	"github.com/go-chi/chi"
)

// Server struct defines the basic layout of a HTTP API server
type Server struct {
	router chi.Router

	// settings and configuration
	basePath       string
	cookieSettings cookies.Settings
	jwtOpts        security.JwtOptions

	//api    api.Login
}

// Create instantiates a new Server instance
func Create(basePath string, config config.AppConfig, version config.VersionInfo) *Server {
	base, err := filepath.Abs(basePath)
	if err != nil {
		panic(fmt.Sprintf("cannot resolve basepath '%s', %v", basePath, err))
	}

	srv := Server{
		basePath: base,
		cookieSettings: cookies.Settings{
			Path:   config.Cookies.Path,
			Domain: config.Cookies.Domain,
			Secure: config.Cookies.Secure,
			Prefix: config.Cookies.Prefix,
		},
		jwtOpts: security.JwtOptions{
			JwtSecret:  config.Sec.JwtSecret,
			JwtIssuer:  config.Sec.JwtIssuer,
			CookieName: config.Sec.CookieName,
			RequiredClaim: security.Claim{
				Name:  config.Sec.Claim.Name,
				URL:   config.Sec.Claim.URL,
				Roles: config.Sec.Claim.Roles,
			},
			RedirectURL:   config.Sec.LoginRedirect,
			CacheDuration: config.Sec.CacheDuration,
		},
		//api:            api.New(base, cookieSettings, version, config.OIDC, config.Sec, repo),
	}
	srv.routes()
	return &srv
}

// Handler returns the used router
func (s *Server) Handler() chi.Router {
	return s.router
}

// --------------------------------------------------------------------------
// internal logic / helpers
// --------------------------------------------------------------------------

func serveStaticDir(r chi.Router, public string, static http.Dir) {
	if strings.ContainsAny(public, "{}*") {
		panic("FileServer does not permit URL parameters.")
	}

	root, _ := filepath.Abs(string(static))
	if _, err := os.Stat(root); os.IsNotExist(err) {
		panic("Static Documents Directory Not Found")
	}

	fs := http.StripPrefix(public, http.FileServer(http.Dir(root)))

	r.Get(public+"*", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		file := strings.Replace(r.RequestURI, public, "", 1)
		// if the file contains URL params, remove everything after ?
		if strings.Contains(file, "?") {
			parts := strings.Split(file, "?")
			if len(parts) == 2 {
				file = parts[0] // use everything before the ?
			}
		}
		if _, err := os.Stat(root + file); os.IsNotExist(err) {
			http.ServeFile(w, r, path.Join(root, "index.html"))
			return
		}
		fs.ServeHTTP(w, r)
	}))
}

func serveStaticFile(r chi.Router, path, filepath string) {
	if path == "" {
		panic("no path for fileServer defined!")
	}
	if strings.ContainsAny(path, "{}*") {
		panic("fileServer does not permit URL parameters.")
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath)
	})

	r.Get(path, handler)
	r.Options(path, handler)
}

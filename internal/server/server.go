// Package server implements the API backend of the bookmark application
package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/bihe/commons-go/cookies"
	"github.com/bihe/commons-go/errors"
	"github.com/bihe/commons-go/handler"
	"github.com/bihe/commons-go/security"

	"github.com/bihe/bookmarks/internal"
	"github.com/bihe/bookmarks/internal/config"
	"github.com/bihe/bookmarks/internal/server/html"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

// Server struct defines the basic layout of a HTTP API server
type Server struct {
	router         chi.Router
	basePath       string
	jwtOpts        security.JwtOptions
	cookieSettings cookies.Settings

	logConfig   config.LogConfig
	environment string

	errorHandler *html.TemplateHandler
	appInfoAPI   *handler.AppInfoHandler
}

// Create instantiates a new Server instance
func Create(basePath string, config config.AppConfig, version internal.VersionInfo, environment string) *Server {
	base, err := filepath.Abs(basePath)
	if err != nil {
		panic(fmt.Sprintf("cannot resolve basepath '%s', %v", basePath, err))
	}

	env := config.Environment
	if environment != "" {
		env = environment
	}

	cookieSettings := cookies.Settings{
		Path:   config.Cookies.Path,
		Domain: config.Cookies.Domain,
		Secure: config.Cookies.Secure,
		Prefix: config.Cookies.Prefix,
	}
	errorReporter := errors.NewReporter(cookieSettings, config.ErrorPath)
	baseHandler := handler.Handler{
		ErrRep: errorReporter,
	}

	appInfo := &handler.AppInfoHandler{
		Handler: baseHandler,
		Version: version.Version,
		Build:   version.Build,
	}
	errHandler := &html.TemplateHandler{
		Handler:        baseHandler,
		Version:        version.Version,
		Build:          version.Build,
		CookieSettings: cookieSettings,
		BasePath:       basePath,
	}

	srv := Server{
		basePath: base,
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
			ErrorPath:     config.ErrorPath,
		},
		cookieSettings: cookieSettings,
		logConfig:      config.Log,
		environment:    env,
		appInfoAPI:     appInfo,
		errorHandler:   errHandler,
	}
	srv.routes()
	return &srv
}

// ServeHTTP turns the server into a http.Handler
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

// use the go-chi logger middleware and redirect request logging to a file
func (s *Server) setupRequestLogging() {

	if s.environment != "Development" {
		var file *os.File
		file, err := os.OpenFile(s.logConfig.RequestPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			panic(fmt.Sprintf("cannot use filepath '%s' as a logfile: %v", s.logConfig.RequestPath, err))
		}
		middleware.DefaultLogger = middleware.RequestLogger(&middleware.DefaultLogFormatter{
			Logger:  log.New(file, "", log.LstdFlags),
			NoColor: true,
		})
	}
}
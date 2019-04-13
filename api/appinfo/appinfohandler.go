package appinfo

import (
	"net/http"

	"github.com/bihe/bookmarks/api"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

// --------------------------------------------------------------------------
// Object definitions
// --------------------------------------------------------------------------

// UserInfo provides data for the frontend about the current user
type UserInfo struct {
	Username      string   `json:"userName"`
	DisplayName   string   `json:"displayName"`
	Roles         []string `json:"roles"`
	Authenticated bool     `json:"authenticated"`
}

// VersionInfo provides information about the application version
type VersionInfo struct {
	Version string `json:"version"`
	Build   string `json:"buildNumber"`
}

// AppInfo holde information about the current user and application version
type AppInfo struct {
	User    UserInfo    `json:"userInfo"`
	Version VersionInfo `json:"versionInfo"`
}

// --------------------------------------------------------------------------
// Framework specific objects
// --------------------------------------------------------------------------

// Response wraps the data structs into a framework response
type Response struct {
	*AppInfo
}

// Render the specific response
func (a Response) Render(w http.ResponseWriter, r *http.Request) error {
	// Pre-processing before a response is marshalled and sent across the wire
	return nil
}

// --------------------------------------------------------------------------
// AppInfo API Routes
// --------------------------------------------------------------------------

// MountRoutes defines the application specific routes
func MountRoutes(version, build string) http.Handler {
	r := chi.NewRouter()
	api := appInfoAPI{version: version, build: build}
	r.Get("/", api.GetAppInfo)

	return r
}

// --------------------------------------------------------------------------
// AppInfo API
// --------------------------------------------------------------------------

type appInfoAPI struct {
	version string
	build   string
}

// GetAppInfo returns information about current user and version of the application
func (a appInfoAPI) GetAppInfo(w http.ResponseWriter, r *http.Request) {
	user := api.User(r)
	userInfo := UserInfo{
		Username:      user.Username,
		DisplayName:   user.DisplayName,
		Roles:         user.Roles,
		Authenticated: true,
	}
	version := "1.0.0"
	build := "localbuild"
	if a.version != "" {
		version = a.version
	}
	if a.build != "" {
		build = a.build
	}
	versionInfo := VersionInfo{
		Version: version,
		Build:   build,
	}
	appInfo := AppInfo{
		User:    userInfo,
		Version: versionInfo,
	}
	render.Render(w, r, Response{AppInfo: &appInfo})
}

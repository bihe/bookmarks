package api

import (
	"net/http"

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

// AppInfoResponse wraps the data structs into a framework response
type AppInfoResponse struct {
	*AppInfo
}

// Render the specific response
func (a AppInfoResponse) Render(w http.ResponseWriter, r *http.Request) error {
	// Pre-processing before a response is marshalled and sent across the wire
	return nil
}

// --------------------------------------------------------------------------
// AppInfo API Routes
// --------------------------------------------------------------------------

// MountRoutes defines the application specific routes
func MountRoutes() http.Handler {
	r := chi.NewRouter()
	r.Get("/", GetAppInfo)

	return r
}

// --------------------------------------------------------------------------
// AppInfo API
// --------------------------------------------------------------------------

// GetAppInfo returns information about current user and version of the application
func GetAppInfo(w http.ResponseWriter, r *http.Request) {
	user := User(r)
	userInfo := UserInfo{
		Username:      user.Username,
		DisplayName:   user.DisplayName,
		Roles:         user.Roles,
		Authenticated: true,
	}
	versionInfo := VersionInfo{
		Version: Version,
		Build:   Build,
	}
	appInfo := AppInfo{
		User:    userInfo,
		Version: versionInfo,
	}
	render.Render(w, r, AppInfoResponse{AppInfo: &appInfo})
}

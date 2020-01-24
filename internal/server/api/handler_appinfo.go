package api

import (
	"fmt"
	"net/http"

	"github.com/bihe/commons-go/security"
	log "github.com/sirupsen/logrus"
)

// swagger:operation GET /appinfo appinfo HandleAppInfo
//
// provides information about the application
//
// meta-data of the application including authenticated user and version
//
// ---
// produces:
// - application/json
// responses:
//   '200':
//     description: Meta
//     schema:
//       "$ref": "#/definitions/Meta"
//   '401':
//     description: ProblemDetail
//     schema:
//       "$ref": "#/definitions/ProblemDetail"
//   '403':
//     description: ProblemDetail
//     schema:
//       "$ref": "#/definitions/ProblemDetail"
func (a *bookmarksAPI) HandleAppInfo(user security.User, w http.ResponseWriter, r *http.Request) error {
	log.WithField("func", "api.HandleAppInfo").Debugf("return the application metadata info")
	info := Meta{
		Version: fmt.Sprintf("%s-%s", a.Version, a.Build),
		UserInfo: UserInfo{
			Email:       user.Email,
			DisplayName: user.DisplayName,
			Roles:       user.Roles,
		},
	}
	a.respond(w, r, http.StatusOK, info)
	return nil
}

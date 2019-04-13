package context

import (
	"net/http"

	"github.com/bihe/bookmarks/core"
	"github.com/bihe/bookmarks/security"
)

// User retrieves the user object from the security middleware
func User(r *http.Request) *security.User {
	user := r.Context().Value(core.ContextUser).(*security.User)
	if user == nil {
		panic("could not get User from context")
	}
	return user
}

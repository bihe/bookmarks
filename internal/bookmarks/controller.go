package bookmarks

import (
	"net/http"

	"github.com/bihe/bookmarks-go/internal/conf"
	"github.com/bihe/bookmarks-go/internal/context"
	"github.com/bihe/bookmarks-go/internal/security"
	"github.com/bihe/bookmarks-go/internal/store"
	"github.com/gin-gonic/gin"
)

// CheckContext is used as a preprocessing step to ensure that the context/environment
// for handlers is setup correctly
func CheckContext() gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, ok := c.Get(conf.ContextUser); !ok {
			context.Abort(c, http.StatusInternalServerError, "Could not get an user from context")
			return
		}

		if _, ok := c.Get(conf.ContextUnitOfWork); !ok {
			context.Abort(c, http.StatusInternalServerError, "Could not get an unitOfWork from context")
			return
		}
	}
}

// Controller combines the API methods of the bookmarks logic
type Controller struct{}

// User returns the authenticated principle of the JWT middleware
func (app *Controller) user(c *gin.Context) *security.User {
	return c.MustGet(conf.ContextUser).(*security.User)
}

// unitOfWork returns the store implementation
func (app *Controller) unitOfWork(c *gin.Context) *store.UnitOfWork {
	return c.MustGet(conf.ContextUnitOfWork).(*store.UnitOfWork)
}

// return an error-message to the client
func (app *Controller) error(c *gin.Context, message string) {
	status := http.StatusInternalServerError
	switch c.NegotiateFormat(gin.MIMEHTML, gin.MIMEJSON, gin.MIMEPlain) {
	case gin.MIMEJSON:
		c.JSON(status, gin.H{
			"status":  status,
			"message": message,
		})
	case gin.MIMEHTML:
		fallthrough
	case gin.MIMEPlain:
		c.String(status, message)
	default:
		c.JSON(status, gin.H{
			"status":  status,
			"message": message,
		})
	}
}

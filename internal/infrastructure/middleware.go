package infrastructure

import (
	"log"
	"net/http"

	"github.com/bihe/bookmarks-go/internal/conf"
	"github.com/bihe/bookmarks-go/internal/context"
	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
)

// RequestID is set for each authenticated request to simplify log correlation
const RequestID string = "HttpRequestTraceId"

// Trace adds a unique ID to a HTTP request
func Trace(logTrace bool) gin.HandlerFunc {

	return func(c *gin.Context) {
		// create a unique request-id used for correlation in the logs
		id := xid.New()
		c.Set(RequestID, id)

		c.Next()

		if logTrace {
			log.Printf("\t[TRACE] request for url %s using id '%s'", c.Request.RequestURI, id.String())
		}
	}
}

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

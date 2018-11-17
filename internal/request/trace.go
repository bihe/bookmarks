package request

import (
	"log"

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

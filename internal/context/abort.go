package context

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Abort stops the processing of the *gin.Context
func Abort(c *gin.Context, status int, message string) {
	switch c.NegotiateFormat(gin.MIMEHTML, gin.MIMEJSON, gin.MIMEPlain) {
	case gin.MIMEJSON:
		c.AbortWithStatusJSON(status, gin.H{
			"status":  status,
			"message": message,
		})
	case gin.MIMEHTML:
		fallthrough
	case gin.MIMEPlain:
		c.String(status, message)
		c.Abort()
	default:
		c.AbortWithStatusJSON(status, gin.H{
			"status":  status,
			"message": message,
		})
	}
}

// AbortAndRedirect stops the processing of the *gin.Context; if the request was text/html
// use the provided redirect url
func AbortAndRedirect(c *gin.Context, status int, message string, redirectURL string) {
	switch c.NegotiateFormat(gin.MIMEHTML, gin.MIMEJSON, gin.MIMEPlain) {
	case gin.MIMEJSON:
		c.AbortWithStatusJSON(status, gin.H{
			"status":  status,
			"message": message,
		})
	case gin.MIMEHTML:
		c.Redirect(http.StatusTemporaryRedirect, redirectURL)
		c.Abort()
	case gin.MIMEPlain:
		c.String(status, message)
		c.Abort()
	default:
		c.AbortWithStatusJSON(status, gin.H{
			"status":  status,
			"message": message,
		})
	}
}

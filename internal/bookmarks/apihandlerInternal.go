package bookmarks

import (
	"fmt"
	"net/http"

	"github.com/bihe/bookmarks-go/internal/store"
	"github.com/gin-gonic/gin"
)

// DebugINIT create some testing bookmarks
func (app *BookmarkController) DebugINIT(c *gin.Context) {
	if c.ClientIP() == "127.0.0.1" {
		for i := 0; i < 10; i++ {
			if err := app.unitOfWork(c).CreateBookmark(store.BookmarkItem{
				DisplayName: fmt.Sprintf("ABC_%d", i),
				Path:        fmt.Sprintf("/%d", i),
				URL:         "http://ab.c.de",
				Type:        store.BookmarkNode,
				SortOrder:   0,
			}); err != nil {
				app.error(c, err.Error())
				return
			}
		}
		c.Status(http.StatusOK)
		return
	}
	c.Status(http.StatusForbidden)
}

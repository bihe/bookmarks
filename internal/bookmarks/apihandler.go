package bookmarks

import (
	"fmt"
	"log"
	"net/http"

	"github.com/bihe/bookmarks-go/internal/store"

	"github.com/gin-gonic/gin"
)

// Bookmark is the view-model returned by the API
type Bookmark struct {
	Path        string `json:"path"`
	DisplayName string `json:"displayName"`
	URL         string `json:"url"`
	NodeID      string `json:"nodeId"`
	SortOrder   uint8  `json:"sortOrder"`
	ItemType    string `json:"itemType"`
}

// DebugInitBookmarks create some testing bookmarks
func (app *Controller) DebugInitBookmarks(c *gin.Context) {
	if c.ClientIP() == "127.0.0.1" {
		for i := 0; i < 10; i++ {
			if err := app.unitOfWork(c).CreateBookmark(&store.BookmarkItem{
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

// GetAllBookmarks retrieves the complete list of bookmarks entries from the store
func (app *Controller) GetAllBookmarks(c *gin.Context) {
	var err error

	log.Printf("Got user from authenticated request: '%s'", app.user(c).Username)

	var bookmarks []store.BookmarkItem
	if bookmarks, err = app.unitOfWork(c).GetAllBookmarks(); err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": err,
		})
		return
	}

	items := mapBookmarks(bookmarks)
	c.JSON(http.StatusOK, items)
}

func mapBookmarks(vs []store.BookmarkItem) []Bookmark {
	vsm := make([]Bookmark, len(vs))
	for i, v := range vs {
		t := ""
		switch v.Type {
		case store.BookmarkFolder:
			t = "folder"
		default:
			t = "node"
		}

		vsm[i] = Bookmark{
			DisplayName: v.DisplayName,
			Path:        v.Path,
			NodeID:      v.ItemID,
			ItemType:    t,
			SortOrder:   v.SortOrder,
			URL:         v.URL,
		}
	}
	return vsm
}

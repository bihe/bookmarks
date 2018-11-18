package bookmarks

import (
	"log"
	"net/http"

	"github.com/bihe/bookmarks-go/internal/security"
	"github.com/bihe/bookmarks-go/internal/store"

	"github.com/gin-gonic/gin"
)

// BookmarkType is a Noder or a Folder
type BookmarkType int

const (
	// Folder for bookmakrs
	Folder BookmarkType = iota
	// Node is a bookmark item
	Node
)

type bookmark struct {
	ID          uint         `json:"id"`
	DisplayName string       `json:"displayName"`
	URL         string       `json:"url"`
	SortOrder   uint8        `json:"sortOrder"`
	ItemType    BookmarkType `json:"itemType"`
}

// DebugInitBookmarks create some testing bookmarks
func DebugInitBookmarks(c *gin.Context) {
	db := c.MustGet("DB").(*store.Database)

	if c.ClientIP() == "127.0.0.1" {
		for i := 0; i < 10; i++ {
			db.DB.Create(&store.BookmarkItem{
				DisplayName: "ABC",
				URL:         "http://ab.c.de",
				Type:        store.BookmarkNode,
				SortOrder:   0,
			})
		}
		c.Status(http.StatusOK)
		return
	}
	c.Status(http.StatusForbidden)
}

// GetAllBookmarks retrieves the complete list of bookmarks entries from the store
func GetAllBookmarks(c *gin.Context) {
	user := c.MustGet("User").(security.User)
	db := c.MustGet("DB").(*store.Database)

	log.Printf("Got user from authenticated request: '%s'", user.Username)

	var bookmarks []store.BookmarkItem
	if result := db.DB.Find(&bookmarks); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": result.Error,
		})
		return
	}

	var items []bookmark
	for _, item := range bookmarks {
		items = append(items, bookmark{
			DisplayName: item.DisplayName,
			ID:          item.ID,
			ItemType:    BookmarkType(item.Type),
			SortOrder:   item.SortOrder,
			URL:         item.URL,
		})
	}

	c.JSON(http.StatusOK, items)
}

package handler

import (
	"log"
	"net/http"

	"github.com/bihe/bookmarks/internal/security"

	"github.com/gin-gonic/gin"
)

type bookmarkType int

const (
	// Folder for bookmakrs
	Folder bookmarkType = iota
	// Node is a bookmark item
	Node
)

type bookmark struct {
	ID          int          `json:"id"`
	DisplayName string       `json:"displayName"`
	URL         string       `json:"url"`
	SortOrder   int          `json:"sortOrder"`
	ItemType    bookmarkType `json:"itemType"`
}

// GetAllBookmarks retrieves the complete list of bookmarks entries from the store
func GetAllBookmarks(c *gin.Context) {
	user := c.MustGet("User").(security.User)

	log.Printf("Got user from authenticated request: '%s'", user.Username)

	data := bookmark{ID: 1, DisplayName: "Bookmark", URL: "http://bookmark.com", SortOrder: 0, ItemType: Node}
	c.JSON(http.StatusOK, data)
}

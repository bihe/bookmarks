package bookmarks

import (
	"fmt"
	"net/http"

	"github.com/bihe/bookmarks-go/internal/conf"
	"github.com/bihe/bookmarks-go/internal/security"
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

// BookmarkController combines the API methods of the bookmarks logic
type BookmarkController struct{}

// User returns the authenticated principle of the JWT middleware
func (app *BookmarkController) user(c *gin.Context) *security.User {
	return c.MustGet(conf.ContextUser).(*security.User)
}

// unitOfWork returns the store implementation
func (app *BookmarkController) unitOfWork(c *gin.Context) *store.UnitOfWork {
	return c.MustGet(conf.ContextUnitOfWork).(*store.UnitOfWork)
}

// return an error-message to the client
func (app *BookmarkController) error(c *gin.Context, message string) {
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

// GetAll retrieves the complete list of bookmarks entries from the store
func (app *BookmarkController) GetAll(c *gin.Context) {
	var err error
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

// Create will save a new bookmark entry
func (app *BookmarkController) Create(c *gin.Context) {
	var bookmark Bookmark
	if err := c.ShouldBindJSON(&bookmark); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": err,
		})
		return
	}
	itemType := store.BookmarkNode
	if bookmark.ItemType == "folder" {
		itemType = store.BookmarkNode
	}
	err := app.unitOfWork(c).CreateBookmark(store.BookmarkItem{
		DisplayName: bookmark.DisplayName,
		Path:        bookmark.Path,
		Type:        itemType,
		URL:         bookmark.URL,
		SortOrder:   bookmark.SortOrder,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": err,
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"status":  http.StatusCreated,
		"message": fmt.Sprintf("bookmark item created: %s/%s", bookmark.Path, bookmark.DisplayName),
	})
}

// GetByID returns a single bookmark item, path param :NodeId
func (app *BookmarkController) GetByID(c *gin.Context) {
	nodeID := c.Param("NodeId")
	if nodeID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "missing parameter NodeId",
		})
		return
	}
	var item *store.BookmarkItem
	var err error
	if item, err = app.unitOfWork(c).GetItemById(nodeID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": err,
		})
		return
	}
	c.JSON(http.StatusOK, mapBookmark(item))
}

func mapBookmark(item *store.BookmarkItem) Bookmark {
	t := ""
	switch item.Type {
	case store.BookmarkFolder:
		t = "folder"
	default:
		t = "node"
	}
	return Bookmark{
		DisplayName: item.DisplayName,
		Path:        item.Path,
		NodeID:      item.ItemID,
		ItemType:    t,
		SortOrder:   item.SortOrder,
		URL:         item.URL,
	}
}

func mapBookmarks(vs []store.BookmarkItem) []Bookmark {
	vsm := make([]Bookmark, len(vs))
	for i, v := range vs {
		vsm[i] = mapBookmark(&v)
	}
	return vsm
}

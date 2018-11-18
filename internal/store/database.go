package store

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// ItemType is used to determine the Item
type ItemType uint8

const (
	// BookmarkNode is a single bookmark
	BookmarkNode ItemType = iota
	// BookmarkFolder is a hierarchy/grouping element
	BookmarkFolder ItemType = iota
)

// BookmarkItem represents an item either, Node or Folder
type BookmarkItem struct {
	gorm.Model
	DisplayName string
	URL         string
	SortOrder   uint8
	Type        ItemType
	Parent      *BookmarkItem
}

// Database wraps the underlying database implementation
type Database struct {
	DB *gorm.DB
}

// OpenConn opens a new DB connection for a request and returns the connection after completion
func OpenConn(connStr string) gin.HandlerFunc {
	return func(c *gin.Context) {
		db, err := gorm.Open("sqlite3", connStr)
		if err != nil {
			abort(c, http.StatusServiceUnavailable, "Could not connect to the database!")
			return
		}
		defer db.Close()
		db.AutoMigrate(&BookmarkItem{})
		c.Set("DB", &Database{DB: db})
		c.Next()
	}
}

func abort(c *gin.Context, status int, message string) {
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

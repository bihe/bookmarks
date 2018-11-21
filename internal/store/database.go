package store

import (
	"net/http"

	"github.com/bihe/bookmarks-go/internal/conf"
	"github.com/bihe/bookmarks-go/internal/context"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/rs/xid"

	// get sqlite db driver
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// ItemType is used to determine the Item
type ItemType uint8

const (
	// BookmarkNode is a single bookmark
	BookmarkNode ItemType = iota
	// BookmarkFolder is a hierarchy/grouping element
	BookmarkFolder ItemType = iota
)

// BookmarkItem represents an item - Node or Folder
// a parent-child structure is modeled
type BookmarkItem struct {
	Path        string `gorm:"primary_key;size:512"`
	ItemID      string `gorm:"unique;not null;size:512"`
	DisplayName string `gorm:"unique;not null;size:128"`
	URL         string `gorm:"not null;size:256"`
	SortOrder   uint8  `gorm:"column:sortorder"`
	Type        ItemType
}

// UnitOfWork wraps the underlying database implementation
type UnitOfWork struct {
	db *gorm.DB
}

// GetAllBookmarks returns all available bookmarks
func (u *UnitOfWork) GetAllBookmarks() ([]BookmarkItem, error) {
	var bookmarks []BookmarkItem
	if result := u.db.Order("path asc").Order("sortorder asc").Find(&bookmarks); result.Error != nil {
		return nil, result.Error
	}
	return bookmarks, nil
}

// CreateBookmark saves a new bookmark in the store
func (u *UnitOfWork) CreateBookmark(item *BookmarkItem) error {
	if item.ItemID == "" {
		item.ItemID = xid.New().String()
	}
	return u.db.Create(&item).Error
}

// InUnitOfWork opens a new DB connection for a request and closes the connection after completion
func InUnitOfWork(connStr string) gin.HandlerFunc {
	return func(c *gin.Context) {
		db, err := gorm.Open("sqlite3", connStr)
		if err != nil {
			context.Abort(c, http.StatusServiceUnavailable, "Could not connect to the database!")
			return
		}
		defer db.Close()
		db.AutoMigrate(&BookmarkItem{})
		c.Set(conf.ContextUnitOfWork, &UnitOfWork{db: db})
		c.Next()
	}
}

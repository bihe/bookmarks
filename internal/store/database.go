package store

import (
	"context"
	"net/http"

	"github.com/bihe/bookmarks-go/internal/conf"
	"github.com/jinzhu/gorm"
	"github.com/rs/xid"

	// get sqlite db driver
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// UnitOfWorkMiddleware encapsulates the db handler creation
type UnitOfWorkMiddleware struct {
	DbDialect string
	ConnStr   string
}

// UnitOfWorkContext wraps the database access
func (u *UnitOfWorkMiddleware) UnitOfWorkContext(next http.Handler) http.Handler {
	db, err := gorm.Open(u.DbDialect, u.ConnStr)
	if err != nil {
		panic("Could not connect to the database!")
	}
	db.AutoMigrate(&BookmarkItem{})
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), conf.ContextUnitOfWork, &UnitOfWork{db: db})
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// ItemType is used to determine the Item
type ItemType uint8

const (
	// BookmarkNode is a single bookmark
	BookmarkNode ItemType = iota
	// BookmarkFolder is a hierarchy/grouping element
	BookmarkFolder
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
func (u *UnitOfWork) CreateBookmark(item BookmarkItem) error {
	if item.ItemID == "" {
		item.ItemID = xid.New().String()
	}
	return u.db.Create(&item).Error
}

// GetItemByID queries the item by the given itemID
func (u *UnitOfWork) GetItemByID(itemID string) (*BookmarkItem, error) {
	var item BookmarkItem
	if err := u.db.Where("item_id = ?", itemID).First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

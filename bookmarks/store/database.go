package store

import (
	"github.com/jinzhu/gorm"
	"github.com/rs/xid"

	// get sqlite db driver
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// NewUnitOfWork create a new instance of the database interaction logic
// by setting up the datbase
func NewUnitOfWork(dbdialect, connstr string) *UnitOfWork {
	db, err := gorm.Open(dbdialect, connstr)
	if err != nil {
		panic("Could not connect to the database!")
	}
	db.AutoMigrate(&BookmarkItem{})
	return &UnitOfWork{db: db}
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

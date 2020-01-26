// Package store is responsible to interact with the storage backend used for bookmarks
// this is done by implementing a repository for the datbase

package store

import "github.com/jinzhu/gorm"

// Repository defines methods to interact with the database
type Repository interface {
	GetAllBookmarks(username string) ([]Bookmark, error)
}

// --------------------------------------------------------------------------
// Repository implementation
// --------------------------------------------------------------------------

// Create a new repository
func Create(db *gorm.DB) Repository {
	return &dbRepository{
		db: db,
	}
}

// compile-time check if the interface is implemented
var _ Repository = (*dbRepository)(nil)

type dbRepository struct {
	db *gorm.DB
}

// GetAllBookmarks retrieves all available bookmarks for the given user
func (r *dbRepository) GetAllBookmarks(username string) ([]Bookmark, error) {
	var bookmarks []Bookmark
	db := r.db.Order("sort_order").Order("display_name").Where(&Bookmark{UserName: username}).Find(&bookmarks)
	return bookmarks, db.Error
}

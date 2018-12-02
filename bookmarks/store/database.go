package store

import (
	"fmt"
	"io/ioutil"

	"github.com/jmoiron/sqlx"
	"github.com/rs/xid"

	// import sqlite driver
	_ "github.com/mattn/go-sqlite3"
)

// UnitOfWork wraps the underlying database implementation
type UnitOfWork struct {
	db *sqlx.DB
}

// NewUnitOfWork create a new instance of the database interaction logic
// by setting up the datbase
func NewUnitOfWork(dbdialect, connstr string) *UnitOfWork {
	db := sqlx.MustConnect(dbdialect, connstr)
	return &UnitOfWork{db: db}
}

// BookmarkItem represents an entry in the database
type BookmarkItem struct {
	ItemID      string `db:"item_id"`
	Path        string `db:"path"`
	DisplayName string `db:"display_name"`
	URL         string `db:"url"`
	SortOrder   uint8  `db:"sort_order"`
}

// GetAllBookmarks returns all available bookmarks
func (u *UnitOfWork) GetAllBookmarks() ([]BookmarkItem, error) {
	var bookmarks []BookmarkItem
	if err := u.db.Select(&bookmarks, "SELECT * FROM bookmark_items ORDER BY path, sort_order ASC"); err != nil {
		return nil, err
	}
	return bookmarks, nil
}

// CreateBookmark saves a new bookmark in the store
func (u *UnitOfWork) CreateBookmark(item BookmarkItem) error {
	if item.ItemID == "" {
		item.ItemID = xid.New().String()
	}
	var err error
	tx := u.db.MustBegin()
	if _, err = tx.NamedExec("INSERT INTO bookmark_items (item_id, path, display_name, url, sort_order) VALUES(:item_id,:path,:display_name,:url,:sort_order)", &BookmarkItem{
		ItemID:      item.ItemID,
		DisplayName: item.DisplayName,
		Path:        item.Path,
		URL:         item.URL,
		SortOrder:   item.SortOrder,
	}); err != nil {
		if txErr := tx.Rollback(); txErr != nil {
			return fmt.Errorf("could not rollback transaction: %v", txErr)
		}
		return fmt.Errorf("could not save bookmark: %v", err)
	}
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("could not commit transaction: %v", err)
	}
	return err
}

// UpdateBookmark overwrites an existing bookmark
func (u *UnitOfWork) UpdateBookmark(item BookmarkItem) error {
	if item.ItemID == "" {
		return fmt.Errorf("no ID for bookmark provided, cannot update")
	}

	var err error
	_, err = u.GetItemByID(item.ItemID)
	if err != nil {
		return fmt.Errorf("could not get bookmark with ID:'%s': %v", item.ItemID, err)
	}
	tx := u.db.MustBegin()
	if _, err = tx.NamedExec("UPDATE bookmark_items SET path=:path,display_name=:display_name,url=:url,sort_order=:sort_order WHERE item_id=:item_id", &BookmarkItem{
		ItemID:      item.ItemID,
		DisplayName: item.DisplayName,
		Path:        item.Path,
		URL:         item.URL,
		SortOrder:   item.SortOrder,
	}); err != nil {
		if err = tx.Rollback(); err != nil {
			return fmt.Errorf("could not rollback transaction: %v", err)
		}
		return fmt.Errorf("could not update bookmark: %v", err)
	}
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("could not commit transaction: %v", err)
	}
	return err
}

// GetItemByID queries the item by the given itemID
func (u *UnitOfWork) GetItemByID(itemID string) (*BookmarkItem, error) {
	var item BookmarkItem
	if err := u.db.Get(&item, "SELECT * FROM bookmark_items WHERE item_id=$1", itemID); err != nil {
		return nil, fmt.Errorf("could not get bookmark by ID:'%s': %v", itemID, err)
	}
	return &item, nil
}

// InitSchema sets the sqlite database schema
func (u *UnitOfWork) InitSchema(ddlFilePath string) error {
	c, err := ioutil.ReadFile(ddlFilePath)
	if err != nil {
		return fmt.Errorf("could not read ddl.sql file: %v", err)
	}
	if _, err := u.db.Exec(string(c)); err != nil {
		return fmt.Errorf("cannot created db schema from file '%s': %v", ddlFilePath, err)
	}
	return nil
}

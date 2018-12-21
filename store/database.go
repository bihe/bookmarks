package store

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/rs/xid"

	// import sqlite driver
	_ "github.com/mattn/go-sqlite3"
)

// UnitOfWork wraps the underlying database implementation
type UnitOfWork struct {
	db *sqlx.DB
}

// New create a new instance of the database interaction logic
// by setting up the datbase
func New(dbdialect, connstr string) *UnitOfWork {
	db := sqlx.MustConnect(dbdialect, connstr)
	return &UnitOfWork{db: db}
}

// ItemType specifies if an entry is a node or a folder
type ItemType uint8

const (
	// Node is a bookmark item with URL
	Node ItemType = iota
	// Folder is a grouping/hierarchy element for bookmarks
	Folder
)

// BookmarkItem represents an entry in the database
type BookmarkItem struct {
	ItemID      string   `db:"item_id"`
	Path        string   `db:"path"`
	DisplayName string   `db:"display_name"`
	URL         string   `db:"url"`
	SortOrder   uint8    `db:"sort_order"`
	Type        ItemType `db:"type"`
	Username    string   `db:"user_name"`
	Created     int32    `db:"created"`
	Modified    int32    `db:"modified"`
}

// --------------------------------------------------------------------------
// READ METHODS
// --------------------------------------------------------------------------

// AllBookmarks returns all available bookmarks
func (u *UnitOfWork) AllBookmarks(username string) ([]BookmarkItem, error) {
	var bookmarks []BookmarkItem
	if username == "" {
		return nil, fmt.Errorf("cannot use empty Username")
	}
	if err := u.db.Select(&bookmarks, "SELECT * FROM bookmark_items WHERE user_name=? ORDER BY path, sort_order ASC", username); err != nil {
		return nil, err
	}
	return bookmarks, nil
}

// BookmarkByID queries the item by the given itemID
func (u *UnitOfWork) BookmarkByID(itemID, username string) (*BookmarkItem, error) {
	var item BookmarkItem
	if itemID == "" {
		return nil, fmt.Errorf("cannot use an empty ID")
	}
	if username == "" {
		return nil, fmt.Errorf("cannot use empty Username")
	}
	if err := u.db.Get(&item, "SELECT * FROM bookmark_items WHERE user_name=? AND item_id=?", username, itemID); err != nil {
		return nil, fmt.Errorf("could not get bookmark by ID:'%s': %v", itemID, err)
	}
	return &item, nil
}

// BookmarkByPath returns the bookmarks located in the given path
// a path is similar to a filesystem path /a/b/c
func (u *UnitOfWork) BookmarkByPath(path, username string) ([]BookmarkItem, error) {
	var bookmarks []BookmarkItem
	if path == "" {
		return nil, fmt.Errorf("cannot use an empty path")
	}
	if username == "" {
		return nil, fmt.Errorf("cannot use empty Username")
	}
	if err := u.db.Select(&bookmarks, "SELECT * FROM bookmark_items WHERE user_name = ? AND path = ? ORDER BY sort_order ASC, display_name ASC", username, path); err != nil {
		return nil, err
	}
	return bookmarks, nil
}

// BookmarkStartsByPath returns the bookmarks located in the given path
// a path is similar to a filesystem path /a/b/c*
func (u *UnitOfWork) BookmarkStartsByPath(path, username string) ([]BookmarkItem, error) {
	var bookmarks []BookmarkItem
	if path == "" {
		return nil, fmt.Errorf("cannot use an empty path")
	}
	if username == "" {
		return nil, fmt.Errorf("cannot use empty Username")
	}
	if err := u.db.Select(&bookmarks, "SELECT * FROM bookmark_items WHERE user_name = ? AND path LIKE ? ORDER BY sort_order ASC, display_name ASC", username, path+"%"); err != nil {
		return nil, err
	}
	return bookmarks, nil
}

// FolderByPathName returns the specified folder "name" located within the given "path"
func (u *UnitOfWork) FolderByPathName(path, name, username string) (*BookmarkItem, error) {
	var item BookmarkItem
	if path == "" {
		return nil, fmt.Errorf("cannot use an empty path")
	}
	if name == "" {
		return nil, fmt.Errorf("cannot use an empty name")
	}
	if username == "" {
		return nil, fmt.Errorf("cannot use empty Username")
	}
	if err := u.db.Get(&item, "SELECT * FROM bookmark_items WHERE user_name = ? AND path=? AND type=? AND display_name=?", username, path, Folder, name); err != nil {
		return nil, fmt.Errorf("could not get bookmark Folder for path '%s' and name '%s': %v", path, name, err)
	}
	return &item, nil
}

// BookmarkByName returns all bookmarks matching the given name
func (u *UnitOfWork) BookmarkByName(name, username string) ([]BookmarkItem, error) {
	var bookmarks []BookmarkItem
	if name == "" {
		return nil, fmt.Errorf("cannot use an empty name")
	}
	if username == "" {
		return nil, fmt.Errorf("cannot use empty Username")
	}
	if err := u.db.Select(&bookmarks, "SELECT * FROM bookmark_items WHERE user_name = ? AND lower(display_name) LIKE ? ORDER BY sort_order ASC, display_name ASC", username, "%"+strings.ToLower(name)+"%"); err != nil {
		return nil, err
	}
	return bookmarks, nil
}

// --------------------------------------------------------------------------
// TRANSACTIONAL METHODS
// --------------------------------------------------------------------------

// CreateBookmark saves a new bookmark in the store
func (u *UnitOfWork) CreateBookmark(item BookmarkItem) (bItem *BookmarkItem, err error) {
	item.ItemID = xid.New().String()
	tx := u.db.MustBegin()
	defer func() {
		switch err {
		case nil:
			err = tx.Commit()
		default:
			log.Printf("could not complete the transaction: %v", err)
			if e := tx.Rollback(); e != nil {
				err = fmt.Errorf("%v; could not rollback transaction: %v", err, e)
			}
		}
	}()

	if _, err = tx.NamedExec("INSERT INTO bookmark_items (item_id, path, display_name, url, sort_order, type, user_name, created) VALUES(:item_id,:path,:display_name,:url,:sort_order,:type,:user_name,:created)", &BookmarkItem{
		ItemID:      item.ItemID,
		DisplayName: item.DisplayName,
		Path:        item.Path,
		URL:         item.URL,
		SortOrder:   item.SortOrder,
		Type:        item.Type,
		Username:    item.Username,
		Created:     int32(time.Now().Unix()),
	}); err != nil {
		err = fmt.Errorf("could not save bookmark: %v", err)
		bItem = nil
		return
	}
	return &item, err
}

// UpdateBookmark overwrites an existing bookmark
func (u *UnitOfWork) UpdateBookmark(item BookmarkItem) (err error) {
	if item.ItemID == "" {
		return fmt.Errorf("no ID for bookmark provided, cannot update")
	}
	_, err = u.BookmarkByID(item.ItemID, item.Username)
	if err != nil {
		return err
	}
	tx := u.db.MustBegin()
	defer func() {
		switch err {
		case nil:
			err = tx.Commit()
		default:
			log.Printf("could not complete the transaction: %v", err)
			if e := tx.Rollback(); e != nil {
				err = fmt.Errorf("%v; could not rollback transaction: %v", err, e)
			}
		}
	}()

	// we will not change the item-type of an existing item - turn a bookmark into a folder? makes no sense
	if _, err = tx.NamedExec("UPDATE bookmark_items SET path=:path,display_name=:display_name,url=:url,sort_order=:sort_order,modified=:modified WHERE item_id=:item_id AND user_name=:user_name", &BookmarkItem{
		ItemID:      item.ItemID,
		Username:    item.Username,
		DisplayName: item.DisplayName,
		Path:        item.Path,
		URL:         item.URL,
		SortOrder:   item.SortOrder,
		Modified:    int32(time.Now().Unix()),
	}); err != nil {
		err = fmt.Errorf("could not update bookmark: %v", err)
		return
	}
	return err
}

// Delete removes a bookmark by the given itemID
func (u *UnitOfWork) Delete(itemID, username string) (err error) {
	if itemID == "" {
		return fmt.Errorf("cannot use empty ID")
	}
	if username == "" {
		return fmt.Errorf("cannot use empty Username")
	}
	tx := u.db.MustBegin()
	defer func() {
		switch err {
		case nil:
			err = tx.Commit()
		default:
			log.Printf("could not complete the transaction: %v", err)
			if e := tx.Rollback(); e != nil {
				err = fmt.Errorf("%v; could not rollback transaction: %v", err, e)
			}
		}
	}()

	var r sql.Result
	if r, err = tx.Exec("DELETE FROM bookmark_items WHERE user_name = ? AND item_id = ?", username, itemID); err != nil {
		err = fmt.Errorf("cannot delete bookmark with ID: %s; error: %v", itemID, err)
		return
	}
	var c int64
	c, err = r.RowsAffected()
	if err != nil {
		log.Printf("Could not delete item '%s': %v", itemID, err)
		err = fmt.Errorf("no items were deleted: %v", err)
		return
	}
	if c == 0 {
		log.Printf("Could not delete item for ID '%s' and Username '%s'", itemID, username)
		err = fmt.Errorf("no items were deleted")
		return
	}
	return nil
}

// DeletePath removes a whole path of items
// it uses the given path and deletes the items which start with the given path
// e.g. /a/b* -- deletes all items with path /a/b, /a/b/c, /a/b/....
func (u *UnitOfWork) DeletePath(path, username string) (err error) {
	if path == "" {
		return fmt.Errorf("cannot use an empty ID")
	}
	if path == "/" {
		return fmt.Errorf("cannot delete the ROOT path '/'")
	}
	if username == "" {
		return fmt.Errorf("cannot use empty Username")
	}

	tx := u.db.MustBegin()
	defer func() {
		switch err {
		case nil:
			err = tx.Commit()
		default:
			log.Printf("could not complete the transaction: %v", err)
			if e := tx.Rollback(); e != nil {
				err = fmt.Errorf("%v; could not rollback transaction: %v", err, e)
			}
		}
	}()

	var r sql.Result
	if r, err = tx.Exec("DELETE FROM bookmark_items WHERE user_name = ? AND path LIKE ?", username, path+"%"); err != nil {
		err = fmt.Errorf("cannot delete items for path: '%s'; error: %v", path, err)
		return
	}

	var c int64
	c, err = r.RowsAffected()
	if err != nil {
		err = fmt.Errorf("no items were deleted for path '%s': %v", path, err)
		return
	}
	if c == 0 {
		err = fmt.Errorf("could not delete item for path '%s' and Username '%s'", path, username)
		return
	}

	// it is also necessary to delete the folder with the given name
	// e.g. if the supplied path is /A/B then we need to delete /A/B*
	// but also delete the item Path /A, DisplayName B, Type Folder
	i := strings.LastIndex(path, "/")
	if i == -1 {
		err = fmt.Errorf("not a valid path, no path seperator '/' found")
		return
	}
	n := path[i+1:]
	if r, err = tx.Exec("DELETE FROM bookmark_items WHERE user_name = ? AND display_name = ? AND type = ?", username, n, Folder); err != nil {
		err = fmt.Errorf("cannot delete item for path: '%s'; error: %v", n, err)
		return
	}

	c, err = r.RowsAffected()
	if err != nil {
		err = fmt.Errorf("cannot delete items for path '%s': %v", n, err)
		return
	}
	if c == 0 {
		log.Printf("Could not delete item for path '%s' and Username '%s'", n, username)
		err = fmt.Errorf("could not delete item for path '%s' and Username '%s'", n, username)
		return
	}
	return nil
}

// --------------------------------------------------------------------------
// HELPER METHODS
// --------------------------------------------------------------------------

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

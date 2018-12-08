package store_test

import (
	"path"
	"path/filepath"
	"testing"

	"github.com/bihe/bookmarks/store"
)

const dbDialect = "sqlite3"
const connStr = ":memory:"

func setupDB() *store.UnitOfWork {
	dir, err := filepath.Abs("../")
	if err != nil {
		panic("Cannot read ddl.sql")
	}
	p := path.Join(dir, "_db/", "ddl.sql")
	uow := store.NewUnitOfWork(dbDialect, connStr)
	uow.InitSchema(p)
	return uow
}

func TestDBBookmarks(t *testing.T) {
	uow := setupDB()
	err := uow.CreateBookmark(store.BookmarkItem{
		DisplayName: "a",
		ItemID:      "id",
		Path:        "/",
		Type:        store.Node,
		URL:         "http://url",
	})
	if err != nil {
		t.Errorf("cannot create bookmark item: %v", err)
	}
	bookmarks, err := uow.AllBookmarks()
	if err != nil {
		t.Errorf("cannot get bookmarks: %v", err)
	}
	if len(bookmarks) < 1 {
		t.Errorf("no bookmarks returned")
	}

	// create a bookmark with same path/displayname
	err = uow.CreateBookmark(store.BookmarkItem{
		DisplayName: "a",
		ItemID:      "id",
		Path:        "/",
		Type:        store.Node,
		URL:         "http://url",
	})
	if err == nil {
		t.Errorf("unique constraint for path/displayname: %v", err)
	}

	// update a given bookmark
	err = uow.UpdateBookmark(store.BookmarkItem{
		DisplayName: "a UPDATE",
		ItemID:      "id",
		Path:        "/",
		Type:        store.Node,
		URL:         "http://url",
	})
	if err != nil {
		t.Errorf("cannot update bookmark: %v", err)
	}

	// verify the update
	var b *store.BookmarkItem
	b, err = uow.BookmarkByID("id")
	if err != nil {
		t.Errorf("cannot get bookmark by id: %v", err)
	}
	if b.DisplayName != "a UPDATE" {
		t.Errorf("the update of the bookmark did not work!")
	}

	// create another bookmark
	err = uow.CreateBookmark(store.BookmarkItem{
		DisplayName: "b",
		ItemID:      "id1",
		Path:        "/path",
		Type:        store.Node,
		URL:         "http://url",
	})
	if err != nil {
		t.Errorf("cannot create bookmark item: %v", err)
	}

	// get bookmarks by path
	var blist []store.BookmarkItem
	blist, err = uow.BookmarkByPath("/")
	if err != nil {
		t.Errorf("cannot get bookmark by path /: %v", err)
	}
	if len(blist) != 1 {
		t.Errorf("1 bookmarks should be returned by path /, got %d", len(blist))
	}

	blist, err = uow.BookmarkByPath("/path")
	if err != nil {
		t.Errorf("cannot get bookmark by path /path: %v", err)
	}
	if len(blist) != 1 {
		t.Errorf("1 bookmarks should be returned by path /path, got %d", len(blist))
	}

	// create a bookmark 'Folder'
	err = uow.CreateBookmark(store.BookmarkItem{
		DisplayName: "a",
		ItemID:      "folder",
		Path:        "/",
		Type:        store.Folder,
	})
	if err != nil {
		t.Errorf("could not create a folder: %v", err)
	}
	b, err = uow.FolderByPathName("/", "a")
	if err != nil {
		t.Errorf("could not find the given folder %s: %v", "/a", err)
	}
	if b.Type != store.Folder {
		t.Errorf("invalid bookmark type returned. expected %d, got %d", store.Folder, b.Type)
	}

}

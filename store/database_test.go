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
	r, err := uow.CreateBookmark(store.BookmarkItem{
		DisplayName: "a",
		Path:        "/",
		Type:        store.Node,
		URL:         "http://url",
		Username:    "A",
	})
	if err != nil {
		t.Errorf("cannot create bookmark item: %v", err)
	}
	bookmarks, err := uow.AllBookmarks("A")
	if err != nil {
		t.Errorf("cannot get bookmarks: %v", err)
	}
	if len(bookmarks) < 1 {
		t.Errorf("no bookmarks returned")
	}

	// create a bookmark with same path/displayname
	_, err = uow.CreateBookmark(store.BookmarkItem{
		DisplayName: "a",
		ItemID:      "id",
		Path:        "/",
		Type:        store.Node,
		URL:         "http://url",
		Username:    "A",
	})
	if err == nil {
		t.Errorf("unique constraint for path/displayname")
	}

	// update a given bookmark
	err = uow.UpdateBookmark(store.BookmarkItem{
		DisplayName: "a UPDATE",
		ItemID:      r.ItemID,
		Path:        "/",
		Type:        store.Node,
		URL:         "http://url",
		Username:    "A",
	})
	if err != nil {
		t.Errorf("cannot update bookmark: %v", err)
	}

	// verify the update
	var b *store.BookmarkItem
	b, err = uow.BookmarkByID(r.ItemID, "A")
	if err != nil {
		t.Errorf("cannot get bookmark by id: %v", err)
	}
	if b.DisplayName != "a UPDATE" {
		t.Errorf("the update of the bookmark did not work!")
	}

	// create another bookmark
	_, err = uow.CreateBookmark(store.BookmarkItem{
		DisplayName: "b",
		Path:        "/path",
		Type:        store.Node,
		URL:         "http://url",
		Username:    "A",
	})
	if err != nil {
		t.Errorf("cannot create bookmark item: %v", err)
	}

	// get bookmarks by path
	var blist []store.BookmarkItem
	blist, err = uow.BookmarkByPath("/", "A")
	if err != nil {
		t.Errorf("cannot get bookmark by path /: %v", err)
	}
	if len(blist) != 1 {
		t.Errorf("1 bookmarks should be returned by path /, got %d", len(blist))
	}

	blist, err = uow.BookmarkByPath("/path", "A")
	if err != nil {
		t.Errorf("cannot get bookmark by path /path: %v", err)
	}
	if len(blist) != 1 {
		t.Errorf("1 bookmarks should be returned by path /path, got %d", len(blist))
	}

	// create a bookmark 'Folder'
	_, err = uow.CreateBookmark(store.BookmarkItem{
		DisplayName: "a",
		ItemID:      "folder",
		Path:        "/",
		Type:        store.Folder,
		Username:    "A",
	})
	if err != nil {
		t.Errorf("could not create a folder: %v", err)
	}
	b, err = uow.FolderByPathName("/", "a", "A")
	if err != nil {
		t.Errorf("could not find the given folder %s: %v", "/a", err)
	}
	if b.Type != store.Folder {
		t.Errorf("invalid bookmark type returned. expected %d, got %d", store.Folder, b.Type)
	}

	// create another bookmark
	r, err = uow.CreateBookmark(store.BookmarkItem{
		DisplayName: "DELETE",
		Path:        "/",
		Type:        store.Node,
		URL:         "http://url",
		Username:    "A",
	})
	if err != nil {
		t.Errorf("cannot create bookmark item: %v", err)
	}

	err = uow.Delete(r.ItemID, "B")
	if err == nil {
		t.Errorf("user 'B' should not be allowed to delete item")
	}

	err = uow.Delete(r.ItemID, "A")
	if err != nil {
		t.Errorf("cannot delete a bookmark item: %v", err)
	}

	// create a hierarchy
	r, err = uow.CreateBookmark(store.BookmarkItem{
		DisplayName: "FOLDER1",
		Path:        "/",
		Type:        store.Folder,
		Username:    "A",
	})
	if err != nil {
		t.Errorf("cannot create bookmark folder: %v", err)
	}
	r, err = uow.CreateBookmark(store.BookmarkItem{
		DisplayName: "FOLDER2",
		Path:        "/FOLDER1",
		Type:        store.Folder,
		Username:    "A",
	})
	if err != nil {
		t.Errorf("cannot create bookmark folder: %v", err)
	}
	r, err = uow.CreateBookmark(store.BookmarkItem{
		DisplayName: "Item1",
		Path:        "/FOLDER1/FOLDER2",
		Type:        store.Node,
		URL:         "http://url.com",
		Username:    "A",
	})

	blist, err = uow.BookmarkByPath("/FOLDER1/FOLDER2", "A")
	if err != nil {
		t.Errorf("could not find the given folder %s: %v", "/a", err)
	}
	if len(blist) != 1 {
		t.Errorf("the result should by 1, got %d", len(blist))
	}
	if err != nil {
		t.Errorf("cannot create bookmark item: %v", err)
	}
	err = uow.DeletePath("/FOLDER1", "A")
	if err != nil {
		t.Errorf("cannot delete bookmark path: %s: %v", "/FOLDER1/FOLDER2", err)
	}

	blist, err = uow.BookmarkByPath("/FOLDER1", "A")
	if err != nil {
		t.Errorf("could not find the given folder %s: %v", "/a", err)
	}
	if len(blist) != 0 {
		t.Errorf("the result should by 0, got %d", len(blist))
	}

	blist, err = uow.BookmarkByPath("/", "A")
	if err != nil {
		t.Errorf("could not find the given folder %s: %v", "/a", err)
	}
	for _, item := range blist {
		if item.Type == store.Folder && item.DisplayName == "FOLDER1" {
			t.Errorf("the path '/FOLDER1' was not fully deleted")
		}
	}

}

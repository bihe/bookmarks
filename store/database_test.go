package store_test

import (
	"path"
	"path/filepath"
	"testing"

	"github.com/bihe/bookmarks/store"
)

const dbDialect = "sqlite3"
const dbConnStr = ":memory:"

func setupDB() *store.UnitOfWork {
	dir, err := filepath.Abs("../")
	if err != nil {
		panic("Cannot read ddl.sql")
	}
	p := path.Join(dir, "_db/", "ddl.sql")
	uow := store.New(dbDialect, dbConnStr)
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
		t.Fatalf("cannot create bookmark item: %v", err)
	}
	bookmarks, err := uow.AllBookmarks("A")
	if err != nil {
		t.Fatalf("cannot get bookmarks: %v", err)
	}
	if len(bookmarks) < 1 {
		t.Fatalf("no bookmarks returned")
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
		t.Fatalf("unique constraint for path/displayname")
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
		t.Fatalf("cannot update bookmark: %v", err)
	}

	// verify the update
	var b *store.BookmarkItem
	b, err = uow.BookmarkByID(r.ItemID, "A")
	if err != nil {
		t.Fatalf("cannot get bookmark by id: %v", err)
	}
	if b.DisplayName != "a UPDATE" {
		t.Fatalf("the update of the bookmark did not work!")
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
		t.Fatalf("cannot create bookmark item: %v", err)
	}

	// get bookmarks by path
	var blist []store.BookmarkItem
	blist, err = uow.BookmarkByPath("/", "A")
	if err != nil {
		t.Fatalf("cannot get bookmark by path /: %v", err)
	}
	if len(blist) != 1 {
		t.Fatalf("1 bookmarks should be returned by path /, got %d", len(blist))
	}

	blist, err = uow.BookmarkByPath("/path", "A")
	if err != nil {
		t.Fatalf("cannot get bookmark by path /path: %v", err)
	}
	if len(blist) != 1 {
		t.Fatalf("1 bookmarks should be returned by path /path, got %d", len(blist))
	}

	// search for bookmarks
	blist, err = uow.BookmarkByName("a", "A")
	if err != nil {
		t.Fatalf("cannot get bookmark by name 'a': %v", err)
	}
	if len(blist) != 1 {
		t.Fatalf("1 bookmarks should be returned by name 'a', got %d", len(blist))
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
		t.Fatalf("could not create a folder: %v", err)
	}
	b, err = uow.FolderByPathName("/", "a", "A")
	if err != nil {
		t.Fatalf("could not find the given folder %s: %v", "/a", err)
	}
	if b.Type != store.Folder {
		t.Fatalf("invalid bookmark type returned. expected %d, got %d", store.Folder, b.Type)
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
		t.Fatalf("cannot create bookmark item: %v", err)
	}

	err = uow.Delete(r.ItemID, "B")
	if err == nil {
		t.Fatalf("user 'B' should not be allowed to delete item")
	}

	err = uow.Delete(r.ItemID, "A")
	if err != nil {
		t.Fatalf("cannot delete a bookmark item: %v", err)
	}

	// create a hierarchy
	r, err = uow.CreateBookmark(store.BookmarkItem{
		DisplayName: "FOLDER1",
		Path:        "/",
		Type:        store.Folder,
		Username:    "A",
	})
	if err != nil {
		t.Fatalf("cannot create bookmark folder: %v", err)
	}
	r, err = uow.CreateBookmark(store.BookmarkItem{
		DisplayName: "FOLDER2",
		Path:        "/FOLDER1",
		Type:        store.Folder,
		Username:    "A",
	})
	if err != nil {
		t.Fatalf("cannot create bookmark folder: %v", err)
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
		t.Fatalf("cannot create bookmark item: %v", err)
	}
	if len(blist) != 1 {
		t.Fatalf("the result should by 1, got %d", len(blist))
	}

	blist, err = uow.BookmarkStartsByPath("/FOLDER1/FOLDER2", "A")
	if err != nil {
		t.Fatalf("cannot get bookmarks by start-path '%s': %v", "/FOLDER1/FOLDER2", err)
	}
	if len(blist) != 1 {
		t.Fatalf("the result should by 1, got %d", len(blist))
	}

	err = uow.DeletePath("/FOLDER1", "A")
	if err != nil {
		t.Fatalf("cannot delete bookmark path: %s: %v", "/FOLDER1/FOLDER2", err)
	}

	blist, err = uow.BookmarkByPath("/FOLDER1", "A")
	if err != nil {
		t.Fatalf("could not find the given folder %s: %v", "/a", err)
	}
	if len(blist) != 0 {
		t.Fatalf("the result should by 0, got %d", len(blist))
	}

	blist, err = uow.BookmarkByPath("/", "A")
	if err != nil {
		t.Fatalf("could not find the given folder %s: %v", "/a", err)
	}
	for _, item := range blist {
		if item.Type == store.Folder && item.DisplayName == "FOLDER1" {
			t.Fatalf("the path '/FOLDER1' was not fully deleted")
		}
	}
}

func TestDBBookmarksPathChildCount(t *testing.T) {
	uow := setupDB()

	// create the following structure
	// /Node1
	// /Folder1
	// /Folder1/Node2
	// /Folder2
	// /Folder1/Folder3
	//
	// expected result
	// /, 3
	// /Folder1, 2

	var err error

	_, err = uow.CreateBookmark(store.BookmarkItem{
		DisplayName: "Node1",
		Path:        "/",
		Type:        store.Node,
		URL:         "http://url",
		Username:    "A",
	})
	if err != nil {
		t.Fatalf("cannot create bookmark item: %v", err)
	}

	_, err = uow.CreateBookmark(store.BookmarkItem{
		DisplayName: "Folder1",
		Path:        "/",
		Type:        store.Folder,
		Username:    "A",
	})
	if err != nil {
		t.Fatalf("cannot create bookmark item: %v", err)
	}

	_, err = uow.CreateBookmark(store.BookmarkItem{
		DisplayName: "Node2",
		Path:        "/Folder1",
		Type:        store.Node,
		URL:         "http://url",
		Username:    "A",
	})
	if err != nil {
		t.Fatalf("cannot create bookmark item: %v", err)
	}

	_, err = uow.CreateBookmark(store.BookmarkItem{
		DisplayName: "Folder2",
		Path:        "/",
		Type:        store.Folder,
		Username:    "A",
	})
	if err != nil {
		t.Fatalf("cannot create bookmark item: %v", err)
	}

	_, err = uow.CreateBookmark(store.BookmarkItem{
		DisplayName: "Folder3",
		Path:        "/Folder1",
		Type:        store.Folder,
		Username:    "A",
	})
	if err != nil {
		t.Fatalf("cannot create bookmark item: %v", err)
	}
	var nc []store.NodeCount
	nc, err = uow.PathChildCount()
	if err != nil {
		t.Fatalf("cannot get the path child count: %v", err)
	}
	if len(nc) == 0 {
		t.Fatalf("no enries returned, expected 2 items!")
	}

	exp := []struct {
		path  string
		count int32
	}{
		{
			path:  "/",
			count: 3,
		},
		{
			path:  "/Folder1",
			count: 2,
		},
	}

	for _, r := range exp {
		for _, c := range nc {
			if c.Path == r.path {
				if c.Count != r.count {
					t.Fatalf("Expected '%d' for path '%s' but got '%d'!", r.count, c.Path, c.Count)
				}
			}
		}
	}

}

package store

import (
	"database/sql"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"

	_ "github.com/jinzhu/gorm/dialects/sqlite" // use sqlite for testing
)

const expectations = "there were unfulfilled expectations: %s"

func mockRepository() (Repository, sqlmock.Sqlmock, error) {
	var (
		db   *sql.DB
		DB   *gorm.DB
		err  error
		mock sqlmock.Sqlmock
	)
	if db, mock, err = sqlmock.New(); err != nil {
		return nil, nil, err
	}
	if DB, err = gorm.Open("mysql", db); err != nil {
		return nil, nil, err
	}
	DB.LogMode(true)
	return Create(DB), mock, nil
}

func repository() (Repository, *gorm.DB, error) {
	var (
		DB  *gorm.DB
		err error
	)
	if DB, err = gorm.Open("sqlite3", ":memory:"); err != nil {
		return nil, nil, err
	}
	// Migrate the schema
	DB.AutoMigrate(&Bookmark{})

	DB.LogMode(true)
	return Create(DB), DB, nil
}

func Test_Mock_GetAllBookmarks(t *testing.T) {
	repo, mock, err := mockRepository()
	if err != nil {
		t.Fatalf("Could not create Repository: %v", err)
	}

	userName := "test"
	now := time.Now().UTC()
	rowDef := []string{"id", "path", "display_name", "url", "sort_order", "type", "user_name", "created", "modified", "child_count", "access_count", "favicon"}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `BOOKMARKS` WHERE")).
		WithArgs(userName).
		WillReturnRows(sqlmock.NewRows(rowDef).
			AddRow("id", "path", "display_name", "url", 0, 0, userName, now, nil, 0, 0, ""))

	bookmarks, err := repo.GetAllBookmarks(userName)
	if err != nil {
		t.Errorf("Could not get bookmarks: %v", err)
	}
	if len(bookmarks) != 1 {
		t.Errorf("Invalid number of bookmarks returned: %d", len(bookmarks))
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf(expectations, err)
	}
}

func TestGetAllBookmarks(t *testing.T) {
	repo, db, err := repository()
	if err != nil {
		t.Fatalf("Could not create Repository: %v", err)
	}
	defer db.Close()

	userName := "test"
	bookmarks, err := repo.GetAllBookmarks(userName)
	if err != nil {
		t.Errorf("Could not get bookmarks: %v", err)
	}
	if len(bookmarks) != 0 {
		t.Errorf("Invalid number of bookmarks returned: %d", len(bookmarks))
	}
}

func TestCreateBookmark(t *testing.T) {
	repo, db, err := repository()
	if err != nil {
		t.Fatalf("Could not create Repository: %v", err)
	}
	defer db.Close()

	userName := "username"
	item := Bookmark{
		DisplayName: "displayName",
		Path:        "/",
		SortOrder:   0,
		Type:        Node,
		URL:         "http://url",
		UserName:    userName,
	}
	bm, err := repo.Create(item)
	if err != nil {
		t.Errorf("Could not create bookmarks: %v", err)
	}

	assert.NotEmpty(t, bm.ID)
	assert.Equal(t, "displayName", bm.DisplayName)
	assert.Equal(t, "http://url", bm.URL)

}

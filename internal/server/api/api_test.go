package api

import "github.com/bihe/bookmarks/internal/store"

// create a mock implementation of the store repository
// tests can use the mock to implement their desired behavior
type mockRepository struct{}

var _ store.Repository = (*mockRepository)(nil)

func (m *mockRepository) InUnitOfWork(fn func(repo store.Repository) error) error {
	return nil
}

func (m *mockRepository) Create(item store.Bookmark) (store.Bookmark, error) {
	return store.Bookmark{}, nil
}

func (m *mockRepository) Update(item store.Bookmark) (store.Bookmark, error) {
	return store.Bookmark{}, nil
}

func (m *mockRepository) Delete(item store.Bookmark) error {
	return nil
}

func (m *mockRepository) DeletePath(path, username string) error {
	return nil
}

func (m *mockRepository) GetAllBookmarks(username string) ([]store.Bookmark, error) {
	return nil, nil
}

func (m *mockRepository) GetBookmarksByPath(path, username string) ([]store.Bookmark, error) {
	return nil, nil
}

func (m *mockRepository) GetBookmarksByPathStart(path, username string) ([]store.Bookmark, error) {
	return nil, nil
}

func (m *mockRepository) GetBookmarksByName(name, username string) ([]store.Bookmark, error) {
	return nil, nil
}

func (m *mockRepository) GetMostRecentBookmarks(username string, limit int) ([]store.Bookmark, error) {
	return nil, nil
}

func (m *mockRepository) GetPathChildCount(path, username string) ([]store.NodeCount, error) {
	return nil, nil
}

func (m *mockRepository) GetBookmarkById(id, username string) (store.Bookmark, error) {
	return store.Bookmark{}, nil
}

func (m *mockRepository) GetFolderByPath(path, username string) (store.Bookmark, error) {
	return store.Bookmark{}, nil
}

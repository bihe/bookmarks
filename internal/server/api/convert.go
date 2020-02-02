package api

import "github.com/bihe/bookmarks/internal/store"

func entityToModel(b store.Bookmark) *Bookmark {
	return &Bookmark{
		ID:          b.ID,
		DisplayName: b.DisplayName,
		Path:        b.Path,
		Type:        entityEnumToModel(b.Type),
		URL:         b.URL,
		SortOrder:   b.SortOrder,
		Created:     b.Created,
		Modified:    b.Modified,
		AccessCount: b.AccessCount,
		ChildCount:  b.ChildCount,
		Favicon:     b.Favicon,
	}
}

func entityEnumToModel(t store.NodeType) NodeType {
	if t == store.Folder {
		return Folder
	}
	return Node
}

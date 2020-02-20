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

func entityListToModel(bms []store.Bookmark) []Bookmark {
	var model = make([]Bookmark, 0)
	for _, b := range bms {
		model = append(model, *entityToModel(b))
	}
	return model
}

func entityEnumToModel(t store.NodeType) NodeType {
	if t == store.Folder {
		return Folder
	}
	return Node
}

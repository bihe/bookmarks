package store

import "time"

type NodeType uint

const (
	Node NodeType = iota
	Folder
)

// Bookmark maps the database table to a struct
type Bookmark struct {
	ID          string     `gorm:"primary_key;Type:varchar(255);column:id"`
	Path        string     `gorm:"Type:varchar(255);column:path;NOT NULL;index:IX_PATH;index:IX_PATH_USER"`
	DisplayName string     `gorm:"Type:varchar(128);column:display_name;NOT NULL"`
	URL         string     `gorm:"Type:varchar(512);column:url;NOT NULL;index:IX_SORT_ORDER"`
	SortOrder   int        `gorm:"column:sort_order;DEFAULT:0;NOT NULL"`
	Type        NodeType   `gorm:"column:type;DEFAULT:0;NOT NULL"`
	UserName    string     `gorm:"Type:varchar(128);column:user_name;NOT NULL;index:IX_USER;index:IX_PATH_USER"`
	Created     time.Time  `gorm:"column:created;NOT NULL"`
	Modified    *time.Time `gorm:"column:modified"`
	ChildCount  int        `gorm:"column:child_count;DEFAULT:0;NOT NULL"`
	AccessCount int        `gorm:"column:access_count;DEFAULT:0;NOT NULL"`
	Favicon     string     `gorm:"Type:varchar(128);column:favicon;NOT NULL"`
}

// TableName specifies the name of the Table used
func (Bookmark) TableName() string {
	return "BOOKMARKS"
}

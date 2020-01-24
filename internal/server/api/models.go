package api

// Meta specifies application metadata
// swagger:model
type Meta struct {
	Version string `json:"version"`
	UserInfo
}

// UserInfo provides information about the currently loged-in user
type UserInfo struct {
	Email       string   `json:"email"`
	DisplayName string   `json:"displayName"`
	Roles       []string `json:"roles"`
}

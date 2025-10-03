package httphandlers

import "time"

type createCommentReq struct {
	ParentID *int64 `json:"parent_id,omitempty"`
	Content  string `json:"content"`
	Author   string `json:"author"`
}

type getCommentsResp struct {
	Comments []comment `json:"comments"`
	Total    int       `json:"total"`
}

type comment struct {
	Content   string     `json:"content"`
	Author    string     `json:"author"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
	Level     int        `json:"level"`
}

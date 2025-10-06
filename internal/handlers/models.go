package httphandlers

import "time"

type createCommentReq struct {
	ParentID *int64 `json:"parent_id,omitempty"`
	Content  string `json:"content"`
	Author   string `json:"author"`
}

type getCommentsReq struct {
	ParentID int64  `form:"parent" binding:"required"`
	Page     int    `form:"page"`
	Limit    int    `form:"limit"`
	Sort     string `form:"sort"`
}

type getCommentsResp struct {
	Comments []comment `json:"comments"`
	Total    int       `json:"total"`
	Page     int       `json:"page"`
	Limit    int       `json:"limit"`
	Pages    int       `json:"pages"`
}

type comment struct {
	Content   string     `json:"content"`
	Author    string     `json:"author"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
	Level     int        `json:"level"`
}

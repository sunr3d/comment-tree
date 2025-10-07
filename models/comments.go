package models

import "time"

type Comment struct {
	ID        int64
	ParentID  *int64
	Content   string
	Author    string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
	Level     int
}

type PagParam struct {
	Page   int
	Limit  int
	Sort   string
	Search string
}

type CommentsRes struct {
	Comments []Comment
	Total    int
	Page     int
	Limit    int
	Pages    int
}

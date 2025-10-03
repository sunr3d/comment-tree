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

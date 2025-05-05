package models

import "time"

type Post struct {
	ID          uint `gorm:"primaryKey"`
	UserID      uint
	User        User
	Description string
	ImageURL    string
	Comments    []Comment
	Likes       []Like
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Comment struct {
	ID        uint `gorm:"primaryKey"`
	PostID    uint
	UserID    uint
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Like struct {
	ID     uint `gorm:"primaryKey"`
	PostID uint
	UserID uint
}

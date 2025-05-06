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
	User      User   `gorm:"foreignKey:UserID"` // Keep this for GORM relations & potential full user data
	UserName  string `gorm:"-"`                 // Add this field for direct name access, ignored by GORM for DB schema
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Like struct {
	ID     uint `gorm:"primaryKey"`
	PostID uint
	UserID uint
}

package models

type Post struct {
	ID          uint `gorm:"primaryKey"`
	UserID      uint
	User        User
	Description string
	ImageURL    string
	Comments    []Comment
	Likes       []Like
}

type Comment struct {
	ID      uint `gorm:"primaryKey"`
	PostID  uint
	UserID  uint
	Content string
}

type Like struct {
	ID     uint `gorm:"primaryKey"`
	PostID uint
	UserID uint
}

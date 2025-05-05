package models

type User struct {
	ID        uint `gorm:"primaryKey"`
	Name      string
	Email     string `gorm:"unique"`
	Provider  string
	Posts     []Post
	SessionID string `gorm:"unique"`
}

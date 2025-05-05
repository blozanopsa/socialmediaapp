package migrations

import (
	"backend/models"

	"gorm.io/gorm"
)

// Seed inserts initial data into the database for development/testing
func Seed(db *gorm.DB) error {
	// Example users
	users := []models.User{
		{Name: "Alice", Email: "alice@example.com", Provider: "microsoft"},
		{Name: "Bob", Email: "bob@example.com", Provider: "microsoft"},
	}
	for _, user := range users {
		db.FirstOrCreate(&user, models.User{Email: user.Email})
	}

	// Example posts
	posts := []models.Post{
		{UserID: 1, Description: "Hello, world!", ImageURL: ""},
		{UserID: 2, Description: "My first post!", ImageURL: ""},
	}
	for _, post := range posts {
		db.FirstOrCreate(&post, models.Post{Description: post.Description})
	}

	return nil
}

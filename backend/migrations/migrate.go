package migrations

import (
	"log"

	"backend/models"

	"gorm.io/gorm"
)

// Migrate runs GORM's AutoMigrate for all models
func Migrate(db *gorm.DB) {
	err := db.AutoMigrate(
		&models.User{},
		&models.Post{},
		&models.Comment{},
		&models.Like{},
	)
	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}
}

package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"backend/configs"
	backendControllers "backend/controllers"
	backendMiddleware "backend/middleware"
	"backend/models"
	backendServices "backend/services"

	"backend/migrations"
	"backend/router"
)

// Import the router from router.go

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}
	// Load config
	cfg := configs.LoadConfig()

	// Build DSN
	dsn := cfg.DBUser + ":" + cfg.DBPassword + "@tcp(" + cfg.DBHost + ":" + cfg.DBPort + ")/" + cfg.DBName + "?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// Auto-migrate models
	err = db.AutoMigrate(&models.User{}, &models.Post{}, &models.Comment{}, &models.Like{})
	if err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	// Automatically run migrations on server start
	migrations.Migrate(db)

	// Initialize services and controllers
	postService := backendServices.NewPostService(db)
	userService := backendServices.NewUserService(db)
	postController := backendControllers.NewPostController(postService)
	userController := backendControllers.NewUserController(userService)
	chatController := backendControllers.NewChatController(db)

	r := router.NewRouter(postController, userController)

	// Public chat endpoints
	r.HandleFunc("/public-chat-messages", chatController.GetMessages).Methods("GET")
	r.HandleFunc("/public-chat-messages", chatController.PostMessage).Methods("POST")

	// Health check
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	log.Println("Starting server on :8080...")
	// Wrap the router with CORS middleware globally
	http.ListenAndServe(":8080", backendMiddleware.CORSMiddleware(r))
}

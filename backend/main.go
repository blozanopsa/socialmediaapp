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

	// Wrap the router with the CORS middleware
	r.Use(backendMiddleware.CORSMiddleware)

	// Public user endpoints (OAuth)
	r.HandleFunc("/auth/login", userController.MicrosoftLogin).Methods("GET")
	r.HandleFunc("/auth/microsoft/callback", userController.MicrosoftCallback).Methods("GET")
	r.HandleFunc("/auth/logout", userController.MicrosoftLogout).Methods("POST")
	r.HandleFunc("/auth/facebook/login", userController.FacebookLogin).Methods("GET")
	r.HandleFunc("/auth/facebook/callback", userController.FacebookCallback).Methods("GET")

	// Post endpoints
	api := r.PathPrefix("/api").Subrouter()
	api.Use(backendMiddleware.AuthMiddleware)
	api.HandleFunc("/posts", postController.CreatePost).Methods("POST")
	api.HandleFunc("/posts", postController.GetAllPosts).Methods("GET")
	api.HandleFunc("/posts/{id}", postController.GetPost).Methods("GET")
	api.HandleFunc("/posts/{id}", postController.UpdatePost).Methods("PUT")
	api.HandleFunc("/posts/{id}", postController.DeletePost).Methods("DELETE")
	api.HandleFunc("/posts/filter", postController.FilterPostsByUser).Methods("GET")

	// Health check
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	log.Println("Starting server on :8080...")
	http.ListenAndServe(":8080", r)
}

package router

import (
	"net/http"

	backendControllers "backend/controllers"
	backendMiddleware "backend/middleware"

	"github.com/gorilla/mux"
)

func NewRouter(postController *backendControllers.PostController, userController *backendControllers.UserController) *mux.Router {
	r := mux.NewRouter()

	// Public user endpoints (OAuth)
	r.HandleFunc("/auth/login", userController.MicrosoftLogin).Methods("GET")
	r.HandleFunc("/auth/microsoft/callback", userController.MicrosoftCallback).Methods("GET")
	r.HandleFunc("/auth/logout", userController.MicrosoftLogout).Methods("POST")
	r.HandleFunc("/auth/facebook/login", userController.FacebookLogin).Methods("GET")
	r.HandleFunc("/auth/facebook/callback", userController.FacebookCallback).Methods("GET")

	// Register the `/auth/session/logout` endpoint for clearing sessionID
	r.HandleFunc("/auth/session/logout", userController.Logout).Methods("POST")

	// User names endpoint for chat bubbles
	r.HandleFunc("/users/names", userController.GetUserNames).Methods("GET")

	// Post endpoints (protected, add real auth middleware later)
	api := r.PathPrefix("/api").Subrouter()
	api.Use(backendMiddleware.AuthMiddleware)
	api.HandleFunc("/posts", postController.CreatePost).Methods("POST")
	api.HandleFunc("/posts", postController.GetAllPosts).Methods("GET")
	api.HandleFunc("/posts/{id}", postController.GetPost).Methods("GET")
	api.HandleFunc("/posts/{id}", postController.UpdatePost).Methods("PUT")
	api.HandleFunc("/posts/{id}", postController.DeletePost).Methods("DELETE")
	api.HandleFunc("/posts/filter", postController.FilterPostsByUser).Methods("GET")
	// Like/unlike endpoints
	api.HandleFunc("/posts/{id}/like", postController.LikePost).Methods("POST")
	api.HandleFunc("/posts/{id}/like", postController.UnlikePost).Methods("DELETE")
	// Comment endpoint
	api.HandleFunc("/posts/{id}/comments", postController.AddComment).Methods("POST")
	// Comment edit/delete endpoints
	api.HandleFunc("/posts/{id}/comments/{commentId}", postController.EditComment).Methods("PUT")
	api.HandleFunc("/posts/{id}/comments/{commentId}", postController.DeleteComment).Methods("DELETE")

	// Exclude `/api/user` from AuthMiddleware
	r.HandleFunc("/api/user", userController.GetUserData).Methods("GET")

	// image kit endpoints
	r.HandleFunc("/imagekit-auth", userController.ImageKitAuth).Methods("GET")

	// Health check
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	return r
}

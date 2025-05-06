package controllers

import (
	"backend/models"
	"backend/services"
	"backend/utils"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type PostController struct {
	Service *services.PostService
}

func NewPostController(service *services.PostService) *PostController {
	return &PostController{Service: service}
}

func (c *PostController) CreatePost(w http.ResponseWriter, r *http.Request) {
	// Support both JSON and multipart form
	contentType := r.Header.Get("Content-Type")
	var post models.Post
	if contentType != "" && (contentType == "application/json" || contentType[:16] == "application/json") {
		// Fallback for old clients
		if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	} else {
		// Parse multipart form
		if err := r.ParseMultipartForm(10 << 20); err != nil { // 10MB max
			http.Error(w, "Could not parse multipart form", http.StatusBadRequest)
			return
		}
		post.Description = r.FormValue("description")
		userIDStr := r.FormValue("userId")
		if userIDStr == "" {
			http.Error(w, "Missing userId", http.StatusBadRequest)
			return
		}
		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			http.Error(w, "Invalid userId", http.StatusBadRequest)
			return
		}
		post.UserID = uint(userID)

		// Handle image file if present
		file, header, err := r.FormFile("image")
		if err == nil && file != nil && header != nil {
			log.Printf("[CreatePost] Received image: %s, size: %d bytes", header.Filename, header.Size)
			relPath, err := utils.SaveUploadedFile(file, header, "./public")
			if err != nil {
				log.Printf("[CreatePost] Failed to save image: %v", err)
				http.Error(w, "Failed to save image", http.StatusInternalServerError)
				return
			}
			log.Printf("[CreatePost] Image saved to: %s", relPath)
			post.ImageURL = relPath
		} else {
			log.Printf("[CreatePost] No image uploaded or error: %v", err)
		}
	}
	if err := c.Service.CreatePost(&post); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(post)
}

func (c *PostController) GetPost(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}
	post, err := c.Service.GetPostByID(uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(post)
}

// GetAllPosts supports ?likedBy={userId} to filter posts liked by a user
func (c *PostController) GetAllPosts(w http.ResponseWriter, r *http.Request) {
	likedBy := r.URL.Query().Get("likedBy")
	if likedBy != "" {
		userID, err := strconv.Atoi(likedBy)
		if err != nil {
			http.Error(w, "Invalid likedBy user ID", http.StatusBadRequest)
			return
		}
		posts, err := c.Service.GetPostsLikedByUser(uint(userID))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(posts)
		return
	}
	posts, err := c.Service.GetAllPosts()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(posts)
}

func (c *PostController) FilterPostsByUser(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("user_id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	posts, err := c.Service.FilterPostsByUser(uint(userID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(posts)
}

func (c *PostController) UpdatePost(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}
	var post models.Post
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	post.ID = uint(id)
	if err := c.Service.UpdatePost(&post); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(post)
}

func (c *PostController) DeletePost(w http.ResponseWriter, r *http.Request) {
	log.Println("DeletePost handler called")
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}
	if err := c.Service.DeletePost(uint(id)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// LikePost handles POST /api/posts/{id}/like
func (c *PostController) LikePost(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}
	var body struct {
		UserID uint `json:"userId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.UserID == 0 {
		http.Error(w, "Missing userId", http.StatusBadRequest)
		return
	}
	if err := c.Service.LikePost(uint(id), body.UserID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// UnlikePost handles DELETE /api/posts/{id}/like
func (c *PostController) UnlikePost(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}
	var body struct {
		UserID uint `json:"userId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.UserID == 0 {
		http.Error(w, "Missing userId", http.StatusBadRequest)
		return
	}
	log.Printf("UnlikePost: postID=%d, userID=%d", id, body.UserID)
	if err := c.Service.UnlikePost(uint(id), body.UserID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// AddComment handles POST /api/posts/{id}/comments
func (c *PostController) AddComment(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}
	var body struct {
		Content string `json:"content"`
		UserID  uint   `json:"userId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.UserID == 0 || body.Content == "" {
		http.Error(w, "Missing userId or content", http.StatusBadRequest)
		return
	}

	// Add the comment - service now populates UserName
	comment, err := c.Service.AddComment(uint(id), body.UserID, body.Content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(comment) // Return the comment object directly
}

// EditComment handles PUT /api/posts/{postId}/comments/{commentId}
func (c *PostController) EditComment(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	commentIdStr := mux.Vars(r)["commentId"]
	postID, err := strconv.Atoi(idStr)
	commentID, err2 := strconv.Atoi(commentIdStr)
	if err != nil || err2 != nil {
		http.Error(w, "Invalid post or comment ID", http.StatusBadRequest)
		return
	}
	var body struct {
		Content string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Content == "" {
		http.Error(w, "Missing content", http.StatusBadRequest)
		return
	}
	comment, err := c.Service.EditComment(uint(postID), uint(commentID), body.Content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(comment)
}

// DeleteComment handles DELETE /api/posts/{postId}/comments/{commentId}
func (c *PostController) DeleteComment(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	commentIdStr := mux.Vars(r)["commentId"]
	postID, err := strconv.Atoi(idStr)
	commentID, err2 := strconv.Atoi(commentIdStr)
	if err != nil || err2 != nil {
		http.Error(w, "Invalid post or comment ID", http.StatusBadRequest)
		return
	}
	if err := c.Service.DeleteComment(uint(postID), uint(commentID)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

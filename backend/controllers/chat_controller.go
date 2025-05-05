package controllers

import (
	"backend/models"
	"encoding/json"
	"net/http"

	"gorm.io/gorm"
)

type ChatController struct {
	DB *gorm.DB
}

func NewChatController(db *gorm.DB) *ChatController {
	return &ChatController{DB: db}
}

// GET /public-chat-messages
func (cc *ChatController) GetMessages(w http.ResponseWriter, r *http.Request) {
	var messages []models.ChatMessage
	cc.DB.Order("id asc").Limit(100).Find(&messages)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}

// POST /public-chat-messages
func (cc *ChatController) PostMessage(w http.ResponseWriter, r *http.Request) {
	var msg models.ChatMessage
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		http.Error(w, "Invalid message", http.StatusBadRequest)
		return
	}
	if msg.User == "" || msg.Text == "" {
		http.Error(w, "Missing user or text", http.StatusBadRequest)
		return
	}
	cc.DB.Create(&msg)
	w.WriteHeader(http.StatusCreated)
}

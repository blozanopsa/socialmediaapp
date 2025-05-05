package models

import "time"

type ChatMessage struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	User      uint      `json:"user"` // now stores user ID
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

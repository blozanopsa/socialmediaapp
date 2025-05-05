package models

type ChatMessage struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	User string `json:"user"`
	Text string `json:"text"`
}

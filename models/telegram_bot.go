package models

import (
	"gorm.io/gorm"
)

// TelegramBot defines the bot data model
type TelegramBot struct {
	gorm.Model
	Token    string `json:"token"`
	RoomID   string `json:"room_id"`
	ThreadID string `json:"thread_id"`
}

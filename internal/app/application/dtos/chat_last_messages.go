package dtos

import (
	"github.com/google/uuid"
)

type ChatLastMessages struct {
	ChatID   uuid.UUID `json:"chat_id"`
	ID       uuid.UUID `json:"id"`
	UserID   uuid.UUID `json:"user_id"`
	Message  string    `json:"message"`
	Timestamp string    `json:"timestamp"`
}

type ChatLastMessagesResponse struct {
	Info []ChatLastMessages `json:"info"`
}
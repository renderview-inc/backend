package entities

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	ChatID    uuid.UUID
	Content   string
	CreatedAt time.Time
}

func NewMessage(id, userID, chatID uuid.UUID, content string, createdAt time.Time) *Message {
	return &Message{
		ID:        id,
		UserID:    userID,
		ChatID:    chatID,
		Content:   content,
		CreatedAt: createdAt,
	}
}

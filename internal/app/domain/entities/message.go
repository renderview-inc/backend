package entities

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	ChatTag   string
	Content   string
	CreatedAt time.Time
}

func NewMessage(id, userID uuid.UUID, chatTag, content string, createdAt time.Time) *Message {
	return &Message{
		ID:        id,
		UserID:    userID,
		ChatTag:   chatTag,
		Content:   content,
		CreatedAt: createdAt,
	}
}

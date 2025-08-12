package dtos

import "github.com/google/uuid"

type Message struct {
	ID      uuid.UUID `json:"id,omitempty"`
	UserID  uuid.UUID  `json:"user_id"`
	ChatID  uuid.UUID  `json:"chat_id"`
	Content string     `json:"content"`
}

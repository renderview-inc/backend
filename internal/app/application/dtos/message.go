package dtos

import "github.com/google/uuid"

type Message struct {
	ID        uuid.UUID `json:"id,omitempty"`
	ReplyToID uuid.UUID `json:"reply_to,omitempty"`
	UserID    uuid.UUID `json:"user_id"`
	ChatTag   string    `json:"chat_tag"`
	Content   string    `json:"content"`
}

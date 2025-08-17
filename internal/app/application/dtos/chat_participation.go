package dtos

import "github.com/google/uuid"

type ChatParticipation struct {
	ChatID uuid.UUID `json:"chat_id"`
	UserID uuid.UUID `json:"user_id"`
}

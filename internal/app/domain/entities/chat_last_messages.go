package entities

import "github.com/google/uuid"

type ChatLastMessages struct {
	ID        uuid.UUID
	ChatID    uuid.UUID
	UserID    uuid.UUID
	Message   string
	Timestamp int64
}

func NewChatLastMessages(id, chatID, userID uuid.UUID, message string, timestamp int64) *ChatLastMessages {
	return &ChatLastMessages{
		ID:        id,
		ChatID:    chatID,
		UserID:    userID,
		Message:   message,
		Timestamp: timestamp,
	}
}
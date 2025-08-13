package entities

import (
	"time"

	"github.com/google/uuid"
)

type Chat struct {
	Tag string
	OwnerId uuid.UUID
	CreatedAt time.Time
	Title string
}

func NewChat(tag string, ownerId uuid.UUID, createdAt time.Time, title string) Chat {
	return Chat{
		Tag: tag,
		OwnerId: ownerId,
		CreatedAt: createdAt,
		Title: title,
	}
}
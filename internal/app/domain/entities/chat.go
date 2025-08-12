package entities

import (
	"time"

	"github.com/google/uuid"
)

type Chat struct {
	Id uuid.UUID
	OwnerId uuid.UUID
	CreatedAt time.Time
	Title string
}

func NewChat(id, ownerId uuid.UUID, createdAt time.Time, title string) Chat {
	return Chat{
		Id: id,
		OwnerId: ownerId,
		CreatedAt: createdAt,
		Title: title,
	}
}
package entities

import (
	"time"

	"github.com/google/uuid"
)

type Chat struct {
	Id        uuid.UUID
	Tag       string
	OwnerId   uuid.UUID
	CreatedAt time.Time
	Title     string
}

func NewChat(tag string, ownerId uuid.UUID, title string) Chat {
	return Chat{
		Id:        uuid.New(),
		Tag:       tag,
		OwnerId:   ownerId,
		CreatedAt: time.Now(),
		Title:     title,
	}
}

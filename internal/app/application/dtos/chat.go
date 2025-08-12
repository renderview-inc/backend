package dtos

import "github.com/google/uuid"

type Chat struct {
	Id      uuid.UUID `json:"id,omitempty"`
	OwnerId uuid.UUID `json:"owner_id"`
	Title   string    `json:"title"`
}
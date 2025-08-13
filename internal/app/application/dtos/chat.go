package dtos

import "github.com/google/uuid"

type Chat struct {
	Tag      string `json:"tag,omitempty"`
	OwnerId uuid.UUID `json:"owner_id"`
	Title   string    `json:"title"`
}
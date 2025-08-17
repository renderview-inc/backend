package dtos

import "github.com/google/uuid"

type ChatRequest struct {
	Tag     string    `json:"tag,omitempty"`
	OwnerId uuid.UUID `json:"owner_id"`
	Title   string    `json:"title"`
}

type ChatResponse struct {
	Id  string `json:"id"`
	Tag string `json:"tag"`
}

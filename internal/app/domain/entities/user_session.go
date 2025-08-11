package entities

import (
	"time"

	"github.com/google/uuid"
)

type UserSession struct {
	Id                   uuid.UUID
	UserID               uuid.UUID
	RefreshTokenHash     string
	CreatedAt            time.Time
	UpdatedAt            time.Time
	RefreshExpiresAt     time.Time
	LastUsedAt           time.Time
	Revoked              bool
	RotatedFromSessionID *uuid.UUID
}

func NewUserSession(id uuid.UUID, userID uuid.UUID, refreshTokenHash string,
	createdAt time.Time, updatedAt time.Time, refreshExpiresAt time.Time,
	lastUsedAt time.Time, revoked bool, rotatedFromSessionID *uuid.UUID) UserSession {
	return UserSession{
		id, userID, refreshTokenHash, createdAt, updatedAt, refreshExpiresAt, lastUsedAt, revoked, rotatedFromSessionID,
	}
}

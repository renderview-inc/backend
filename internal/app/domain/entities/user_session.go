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

func (us *UserSession) GetID() uuid.UUID {
	return us.Id
}

func (us *UserSession) GetUserID() uuid.UUID {
	return us.UserID
}

func (us *UserSession) GetRefreshTokenHash() string {
	return us.RefreshTokenHash
}

func (us *UserSession) GetCreatedAt() time.Time {
	return us.CreatedAt
}

func (us *UserSession) GetUpdatedAt() time.Time {
	return us.UpdatedAt
}

func (us *UserSession) GetRefreshExpiresAt() time.Time {
	return us.RefreshExpiresAt
}

func (us *UserSession) GetLastUsedAt() time.Time {
	return us.LastUsedAt
}

func (us *UserSession) GetRevoked() bool {
	return us.Revoked
}

func (us *UserSession) GetRotatedFromSessionGetID() *uuid.UUID {
	return us.RotatedFromSessionID
}

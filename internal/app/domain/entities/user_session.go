package entities

import (
	"time"

	"github.com/google/uuid"
)

type UserSession struct {
	id                   uuid.UUID
	userID               uuid.UUID
	refreshTokenHash     string
	createdAt            time.Time
	updatedAt            time.Time
	refreshExpiresAt     time.Time
	lastUsedAt           time.Time
	revoked              bool
	rotatedFromSessionID *uuid.UUID
}

func NewUserSession(id uuid.UUID, userID uuid.UUID, refreshTokenHash string,
	createdAt time.Time, updatedAt time.Time, refreshExpiresAt time.Time,
	lastUsedAt time.Time, revoked bool, rotatedFromSessionID *uuid.UUID) UserSession {
	return UserSession{
		id, userID, refreshTokenHash, createdAt, updatedAt, refreshExpiresAt, lastUsedAt, revoked, rotatedFromSessionID,
	}
}

func (us *UserSession) ID() uuid.UUID {
	return us.id
}

func (us *UserSession) UserID() uuid.UUID {
	return us.userID
}

func (us *UserSession) RefreshTokenHash() string {
	return us.refreshTokenHash
}

func (us *UserSession) CreatedAt() time.Time {
	return us.createdAt
}

func (us *UserSession) UpdatedAt() time.Time {
	return us.updatedAt
}

func (us *UserSession) RefreshExpiresAt() time.Time {
	return us.refreshExpiresAt
}

func (us *UserSession) LastUsedAt() time.Time {
	return us.lastUsedAt
}

func (us *UserSession) Revoked() bool {
	return us.revoked
}

func (us *UserSession) RotatedFromSessionID() *uuid.UUID {
	return us.rotatedFromSessionID
}

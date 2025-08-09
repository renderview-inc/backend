package entities

import "github.com/google/uuid"

type UserAccount struct {
	id           uuid.UUID
	tag          string
	name         string
	desc         string
	passwordHash string
	email        string
	phone        string
}

func NewUserAccount(id uuid.UUID, tag string, name string, desc string,
	passwordHash string, email string, phone string) *UserAccount {
	return &UserAccount{
		id, tag, name, desc, passwordHash, email, phone,
	}
}

func (uacc *UserAccount) ID() uuid.UUID {
	return uacc.id
}

func (uacc *UserAccount) Tag() string {
	return uacc.tag
}

func (uacc *UserAccount) Name() string {
	return uacc.name
}

func (uacc *UserAccount) Desc() string {
	return uacc.desc
}

func (uacc *UserAccount) PasswordHash() string {
	return uacc.passwordHash
}

func (uacc *UserAccount) Email() string {
	return uacc.email
}

func (uacc *UserAccount) Phone() string {
	return uacc.phone
}

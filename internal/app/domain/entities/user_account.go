package entities

import "github.com/google/uuid"

type UserAccount struct {
	Id           uuid.UUID
	Tag          string
	Name         string
	Desc         string
	PasswordHash string
	Email        string
	Phone        string
}

func NewUserAccount(id uuid.UUID, tag string, name string, desc string,
	passwordHash string, email string, phone string) *UserAccount {
	return &UserAccount{
		id, tag, name, desc, passwordHash, email, phone,
	}
}
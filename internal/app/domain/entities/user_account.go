package entities

import "github.com/google/uuid"

type UserAccount struct {
	Id           uuid.UUID
	Tag         string
	Name        string
	Desc        string
	PasswordHash string
	Email       string
	Phone       string
}

func NewUserAccount(id uuid.UUID, tag string, name string, desc string,
	passwordHash string, email string, phone string) *UserAccount {
	return &UserAccount{
		id, tag, name, desc, passwordHash, email, phone,
	}
}

func (uacc *UserAccount) GetID() uuid.UUID {
	return uacc.Id
}

func (uacc *UserAccount) GetTag() string {
	return uacc.Tag
}

func (uacc *UserAccount) GetName() string {
	return uacc.Name
}

func (uacc *UserAccount) GetDesc() string {
	return uacc.Desc
}

func (uacc *UserAccount) GetPasswordHash() string {
	return uacc.PasswordHash
}

func (uacc *UserAccount) GetEmail() string {
	return uacc.Email
}

func (uacc *UserAccount) GetPhone() string {
	return uacc.Phone
}

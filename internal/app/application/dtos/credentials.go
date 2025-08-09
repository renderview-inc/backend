package dtos

type CredentialsDto struct {
	Email    string `json:"email" validate:"omitempty,email"`
	Phone    string `json:"phone" validate:"omitempty,e164"`
	Tag      string `json:"tag" validate:"omitempty,min=3"`
	Password string `json:"password" validate:"required,min=8"`
}

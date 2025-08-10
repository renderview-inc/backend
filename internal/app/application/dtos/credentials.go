package dtos

type CredentialsDto struct {
	Email    string `json:"email" validate:"omitempty,required_without=Phone,email"`
	Phone    string `json:"phone" validate:"omitempty,required_without=Email,e164"`
	Tag      string `json:"tag" validate:"required,min=3,max=12,regexp=^[a-zA-Z0-9_.]*$"`
	Password string `json:"password" validate:"required,min=8"`
}

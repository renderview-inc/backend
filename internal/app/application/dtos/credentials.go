package dtos

type Credentials struct {
	Email    string `json:"email" validate:"omitempty,required_without=Phone,email"`
	Phone    string `json:"phone" validate:"omitempty,required_without=Email,e164"`
	Tag      string `json:"tag" validate:"required,matches"`
	Password string `json:"password" validate:"required,min=8"`
}

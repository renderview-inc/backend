package dtos

type RegisterDto struct {
	Credentials CredentialsDto `json:"credentials"`
	Name        string         `json:"name"`
	Desc        string         `json:"description"`
}

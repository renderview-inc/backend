package dtos

type Register struct {
	Credentials Credentials    `json:"credentials"`
	Name        string         `json:"name"`
	Desc        string         `json:"description"`
}

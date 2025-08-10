package dtos

type LoginDto struct {
	Credentials CredentialsDto `json:"credentials"`
	LoginMeta   LoginMetaDto   `json:"login_metadata"`
}

package services

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"
)

type Base64TokenIssuer struct {
	tokenLen             int32
	accessTokenLifeTime  time.Duration
	refreshTokenLifeTime time.Duration
}

func NewBase64TokenIssuer(tokenLen int32, accessTokenLifeTime time.Duration, refreshTokenLifeTime time.Duration) Base64TokenIssuer {
	return Base64TokenIssuer{tokenLen, accessTokenLifeTime, refreshTokenLifeTime}
}

func (ti Base64TokenIssuer) issueToken() (string, error) {
	bytes := make([]byte, ti.tokenLen)
	_, err := rand.Read(bytes)

	if err != nil {
		return "", err
	}

	encodedToken := base64.RawURLEncoding.EncodeToString(bytes)
	return encodedToken, nil
}

func (ti Base64TokenIssuer) IssueAccessToken() (string, time.Duration, error) {
	token, err := ti.issueToken()
	if err != nil {
		return "", time.Duration(0), err
	}

	return token, ti.accessTokenLifeTime, nil
}

func (ti Base64TokenIssuer) IssueRefreshToken(sessionID string) (string, time.Duration, error) {
	token, err := ti.issueToken()
	if err != nil {
		return "", time.Duration(0), err
	}

	fullToken := fmt.Sprintf("%s.%s", sessionID, token)

	return fullToken, ti.refreshTokenLifeTime, nil
}

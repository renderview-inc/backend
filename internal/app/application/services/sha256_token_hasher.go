package services

import (
	"crypto/sha256"
	"encoding/hex"
)

type Sha256TokenHasher struct{}

func NewSha256TokenHasher() Sha256TokenHasher {
	return Sha256TokenHasher{}
}

func (sth Sha256TokenHasher) HashToken(token string) (string, error) {
	hash := sha256.Sum256([]byte(token))
	hexHash := hex.EncodeToString(hash[:])

	return hexHash, nil
}

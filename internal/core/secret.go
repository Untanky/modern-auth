package core

import (
	"encoding/base64"

	"github.com/Untanky/modern-auth/internal/utils"
)

type SecretValue struct {
	value []byte
}

func NewSecretValue(value string) *SecretValue {
	return &SecretValue{value: []byte(value)}
}

func (s *SecretValue) MarshalJSON() ([]byte, error) {
	return []byte("\"" + string(s.String()) + "\""), nil
}

func (s *SecretValue) String() string {
	hash := utils.HashShake256(s.value)
	base64Encoded := base64.StdEncoding.EncodeToString(hash)
	return base64Encoded
}

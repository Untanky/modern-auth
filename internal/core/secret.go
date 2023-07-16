package core

import (
	"encoding/base64"

	"golang.org/x/crypto/sha3"
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
	hash := hash_fast(s.value)
	base64Encoded := base64.StdEncoding.EncodeToString(hash)
	return base64Encoded
}

func hash_fast(s []byte) []byte {
	h := make([]byte, 64)
	sha3.ShakeSum256(h, s)
	return h
}

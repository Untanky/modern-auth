package utils

import (
	"crypto"

	"golang.org/x/crypto/sha3"
)

// function to hash a string
func HashSHA256(s []byte) []byte {
	hash := crypto.SHA256.New()
	hash.Write(s)
	return hash.Sum(nil)
}

func HashShake256(value []byte) []byte {
	hash := make([]byte, 64)
	sha3.ShakeSum256(hash, value)
	return hash
}

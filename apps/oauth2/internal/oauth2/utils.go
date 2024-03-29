package oauth2

import (
	"crypto"
	"math/rand"
)

// function to generate a random string
//
// Deprecated: use utils.RandomString instead
func randomString(size int) string {
	const randomStringChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, size)
	for i := range b {
		b[i] = randomStringChars[rand.Intn(len(randomStringChars))]
	}
	return string(b)
}

// function to hash a string
//
// Deprecated: use utils.HashSHA256 instead
func hash(s string) string {
	hash := crypto.SHA256.New()
	hash.Write([]byte(s))
	return string(hash.Sum(nil))
}

package utils

import "math/rand"

var alphabet = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func RandomBytes(bytes []byte) {
	for i := range bytes {
		bytes[i] = alphabet[rand.Intn(len(alphabet))]
	}
}

// function to generate a random string
func RandomString(size int) string {
	data := make([]byte, size)
	RandomBytes(data)
	return string(data)
}

package generator

import (
	"crypto/rand"
	"encoding/base64"
)

func RandomPassword(length int) (string, error) {
	// Calculate the number of bytes needed
	byteLength := length * 6 / 8
	if length%8 != 0 {
		byteLength++
	}

	// Generate random bytes
	randomBytes := make([]byte, byteLength)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	// Encode to base64
	password := base64.RawURLEncoding.EncodeToString(randomBytes)

	// Trim to desired length
	return password[:length], nil
}

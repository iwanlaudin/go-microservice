package helpers

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net"
	"net/mail"
	"strings"
)

// HashString generates a base64 encoded SHA-256 hash of a string.
func HashString(input string) string {
	hash := sha256.Sum256([]byte(input))
	return base64.StdEncoding.EncodeToString(hash[:])
}

// GenerateRandomBytes generates random bytes of a given size.
func GenerateRandomBytes(size int) ([]byte, error) {
	b := make([]byte, size)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// GenerateRandomString generates a random string of a given size.
func GenerateRandomString(size int) (string, error) {
	bytes, err := GenerateRandomBytes(size)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(bytes), nil
}

// Validate email
func ValidateMail(email string) error {
	_, err := mail.ParseAddress(email)
	if err != nil {
		return fmt.Errorf("invalid email")
	}

	parts := strings.Split(email, "@")

	_, err = net.LookupMX(parts[1])
	if err != nil {
		return fmt.Errorf("domain not found")
	}

	return nil
}

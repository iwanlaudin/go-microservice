package helpers

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"net"
	"net/mail"
	"strings"

	"github.com/google/uuid"
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
		return CustomError("Invalid email")
	}

	parts := strings.Split(email, "@")

	_, err = net.LookupMX(parts[1])
	if err != nil {
		return CustomError("Domain not found")
	}

	return nil
}

func ConvertUserIDToUUID(userId string) (uuid.UUID, error) {
	id, err := uuid.Parse(userId)
	if err != nil {
		return uuid.Nil, CustomError("Failed to convert userId to UUID")
	}
	return id, nil
}

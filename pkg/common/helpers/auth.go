package helpers

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"
)

func GeneratePasswordHash(salt string, password string) string {
	combined := salt + password

	hash := sha512.Sum512([]byte(combined))

	return base64.StdEncoding.EncodeToString(hash[:])
}

func VerifyPasswordHash(passwordhash string, password string, salt string) bool {
	hash := GeneratePasswordHash(salt, password)

	return hash == passwordhash
}

func GenerateSaltHash(size int) (string, error) {
	salt := make([]byte, size)

	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(salt), nil
}

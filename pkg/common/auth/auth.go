package auth

import (
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       string
	Username string
	Role     string
}

func Authenticate(username, password string) (*User, error) {
	if username == "admin" && checkPasswordHash(password, "$2a$10$...") {
		return &User{ID: "1", Username: "admin", Role: "admin"}, nil
	}
	return nil, ErrInvalidCredentials
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GeneratePasswordHash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

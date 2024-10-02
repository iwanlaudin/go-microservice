package auth

import "errors"

var (
	ErrInvalidCredentials     = errors.New("invalid credentials")
	ErrInvalidToken           = errors.New("invalid token")
	ErrExpiredToken           = errors.New("expired token")
	ErrInsufficientPermission = errors.New("insufficient permission")
)

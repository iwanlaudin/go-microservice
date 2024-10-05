package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/iwanlaudin/go-microservice/pkg/common/config"
	"github.com/iwanlaudin/go-microservice/pkg/common/helpers"
)

var (
	ErrInvalidToken           = errors.New("invalid token")
	ErrExpiredToken           = errors.New("expired token")
	ErrInvalidExpiredClaim    = errors.New("invalid expiration claim")
	ErrInvalidRoleClaim       = errors.New("invalid role claim")
	ErrInvalidUsernameClaim   = errors.New("invalid username claim")
	ErrInvalidIdClaim         = errors.New("invalid user_id claim")
	ErrInsufficientPermission = errors.New("insufficient permission")
)

func GenerateToken(userClaim map[string]interface{}) (string, error) {
	if userClaim["id"] == "" || userClaim["username"] == "" || userClaim["email"] == "" {
		return "", helpers.CustomError("invalid user data")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       userClaim["id"],
		"username": userClaim["username"],
		"email":    userClaim["email"],
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
		"iat":      time.Now().Unix(),
	})

	secretKey, err := config.GetSecretKey()
	if err != nil {
		return "", helpers.CustomError("failed to get secret key: %w", err)
	}

	signedToken, err := token.SignedString(secretKey)
	if err != nil {
		return "", helpers.CustomError("failed to sign token: %w", err)
	}

	return signedToken, nil
}

func validateClaims(token *jwt.Token) (map[string]interface{}, error) {
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		id, ok := claims["id"].(string)
		if !ok {
			return nil, ErrInvalidIdClaim
		}

		username, ok := claims["username"].(string)
		if !ok {
			return nil, ErrInvalidUsernameClaim
		}

		email, ok := claims["email"].(string)
		if !ok {
			return nil, ErrInvalidRoleClaim
		}

		exp, ok := claims["exp"].(float64)
		if !ok {
			return nil, ErrInvalidExpiredClaim
		}

		if time.Now().Unix() > int64(exp) {
			return nil, ErrExpiredToken
		}

		userContext := map[string]interface{}{
			"id":       id,
			"username": username,
			"email":    email,
		}

		return userContext, nil
	}

	return nil, ErrInvalidToken
}

func ValidateToken(tokenString string) (map[string]interface{}, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, helpers.CustomError("unexpected signing method: %v", token.Header["alg"])
		}

		jwtKey, err := config.GetSecretKey()
		if err != nil {
			return nil, helpers.CustomError("failed to get secret key: %w", err)
		}
		return jwtKey, nil
	})

	if err != nil {
		return nil, err
	}

	if user, err := validateClaims(token); err != nil {
		return nil, err
	} else {
		return user, nil
	}
}

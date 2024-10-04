package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/iwanlaudin/go-microservice/pkg/common/config"
)

func GenerateToken(user *User) (string, error) {
	if user.ID == "" || user.Username == "" || user.Role == "" {
		return "", fmt.Errorf("invalid user data")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"role":     user.Role,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
		"iat":      time.Now().Unix(),
	})

	secretKey, err := config.GetSecretKey()
	if err != nil {
		return "", fmt.Errorf("failed to get secret key: %w", err)
	}

	signedToken, err := token.SignedString(secretKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return signedToken, nil
}

func ValidateToken(tokenString string) (*User, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		jwtKey, err := config.GetSecretKey()
		if err != nil {
			return nil, fmt.Errorf("failed to get secret key: %w", err)
		}
		return jwtKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userId, ok := claims["user_id"].(string)
		if !ok {
			return nil, fmt.Errorf("invalid username claim")
		}

		username, ok := claims["username"].(string)
		if !ok {
			return nil, fmt.Errorf("invalid username claim")
		}

		role, ok := claims["role"].(string)
		if !ok {
			return nil, fmt.Errorf("invalid role claim")
		}

		exp, ok := claims["exp"].(float64)
		if !ok {
			return nil, fmt.Errorf("invalid expiration claim")
		}

		if time.Now().Unix() > int64(exp) {
			return nil, fmt.Errorf("token has expired")
		}

		user := &User{
			ID:       userId,
			Username: username,
			Role:     role,
		}

		return user, nil
	}

	return nil, ErrInvalidToken
}

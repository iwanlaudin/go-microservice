package response

import "github.com/google/uuid"

type UserTokenResponse struct {
	UserId       uuid.UUID `json:"user_id"`
	Expiry       float64   `json:"expiry"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
}

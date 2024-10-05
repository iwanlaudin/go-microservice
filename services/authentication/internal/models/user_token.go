package models

import (
	"time"

	"github.com/google/uuid"
)

type UserToken struct {
	ID           uuid.UUID `db:"id" json:"id"`
	UserId       uuid.UUID `db:"user_id" json:"user_id"`
	RefreshToken string    `db:"refresh_token" json:"refresh_token"`
	ExpiryAt     time.Time `db:"expiry_at" json:"expiry_at"`
	UsedAt       time.Time `db:"used_at" json:"used_at"`
	IsUsed       bool      `db:"is_used" json:"is_used"`
	CreatedBy    *string   `db:"created_by" json:"created_by,omitempty"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedBy    *string   `db:"updated_by" json:"updated_by,omitempty"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
	IsDeleted    bool      `db:"is_deleted" json:"is_deleted"`
}

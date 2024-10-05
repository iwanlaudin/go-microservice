package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `db:"id" json:"id"`
	FirstName string    `db:"first_name" json:"first_name"`
	LastName  string    `db:"last_name,omitempty" json:"last_name,omitempty"`
	Username  string    `db:"username" json:"username"`
	Email     string    `db:"email" json:"email"`
	Salt      string    `db:"salt" json:"-"`
	Password  string    `db:"password" json:"-"`
	CreatedBy *string   `db:"created_by" json:"created_by,omitempty"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedBy *string   `db:"updated_by" json:"updated_by,omitempty"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
	IsDeleted bool      `db:"is_deleted" json:"is_deleted"`

	// Nested UserToken struct
	UserToken UserToken `json:"user_token,omitempty"`
}

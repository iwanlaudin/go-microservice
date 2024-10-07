package models

import (
	"time"

	"github.com/google/uuid"
)

type Ticket struct {
	ID        uuid.UUID  `db:"id" json:"ticket_id"`
	EventID   uuid.UUID  `db:"event_id" json:"event_id"`
	UserID    uuid.UUID  `db:"user_id" json:"user_id"`
	Quantity  int        `db:"quantity" json:"quantity"`
	Status    string     `db:"status" json:"status"`
	CreatedBy *string    `db:"created_by" json:"created_by,omitempty"`
	CreatedAt *time.Time `db:"created_at" json:"created_at"`
	UpdatedBy *string    `db:"updated_by" json:"updated_by,omitempty"`
	UpdatedAt *time.Time `db:"updated_at" json:"updated_at"`
	IsDeleted bool       `db:"is_deleted" json:"is_deleted"`
}

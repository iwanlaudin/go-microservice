package models

import (
	"time"

	"github.com/google/uuid"
)

type Event struct {
	ID               uuid.UUID  `db:"id" json:"event_id"`
	Name             string     `db:"name" json:"event_name"`
	Date             time.Time  `db:"date" json:"event_date"`
	Location         string     `db:"location" json:"location"`
	AvailableTickets int        `db:"available_tickets" json:"available_tickets"`
	CreatedBy        *string    `db:"created_by" json:"created_by,omitempty"`
	CreatedAt        *time.Time `db:"created_at" json:"created_at"`
	UpdatedBy        *string    `db:"updated_by" json:"updated_by,omitempty"`
	UpdatedAt        *time.Time `db:"updated_at" json:"updated_at"`
	IsDeleted        bool       `db:"is_deleted" json:"is_deleted"`
}

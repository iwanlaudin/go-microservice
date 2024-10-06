package response

import (
	"time"

	"github.com/google/uuid"
)

type EventResponse struct {
	ID               uuid.UUID  `json:"event_id"`
	Name             string     `json:"event_name"`
	Date             time.Time  `json:"event_date"`
	Location         string     `json:"location"`
	AvailableTickets int        `json:"available_tickets"`
	CreatedBy        *string    `json:"created_by,omitempty"`
	CreatedAt        *time.Time `json:"created_at,omitempty"`
	UpdatedBy        *string    `json:"updated_by,omitempty"`
	UpdatedAt        *time.Time `json:"updated_at,omitempty"`
}

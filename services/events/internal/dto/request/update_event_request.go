package request

import (
	"time"

	"github.com/google/uuid"
)

type UpdateEventRequest struct {
	ID               uuid.UUID `json:"event_id"`
	Name             string    `json:"event_name"`
	Date             time.Time `json:"event_date"`
	Location         string    `json:"location"`
	AvailableTickets int       `json:"available_tickets"`
}

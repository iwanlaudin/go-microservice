package request

import (
	"time"
)

type CreateEventRequest struct {
	Name             string    `json:"event_name"`
	Date             time.Time `json:"event_date"`
	Location         string    `json:"location"`
	AvailableTickets int       `json:"available_tickets"`
}

package publisher

import "github.com/google/uuid"

type TicketReservedEvent struct {
	ID       uuid.UUID `json:"ticket_id"`
	EventID  uuid.UUID `json:"event_id"`
	UserID   uuid.UUID `json:"user_id"`
	Quantity int       `json:"quantity"`
	Status   string    `json:"status"`
}

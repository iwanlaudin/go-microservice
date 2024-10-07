package messaging

import "github.com/google/uuid"

type TicketReservedMessage struct {
	ID       uuid.UUID `json:"ticket_id"`
	EventID  uuid.UUID `json:"event_id"`
	UserID   uuid.UUID `json:"user_id"`
	Quantity int       `json:"quantity"`
	Status   string    `json:"status"`
}

package response

import (
	"time"

	"github.com/google/uuid"
)

type TicketResponse struct {
	ID        uuid.UUID  `json:"ticket_id"`
	EventID   uuid.UUID  `json:"event_id"`
	UserID    uuid.UUID  `json:"user_id"`
	Quantity  int        `json:"quantity"`
	Status    string     `json:"status"`
	CreatedBy *string    `json:"created_by,omitempty"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedBy *string    `json:"updated_by,omitempty"`
	UpdatedAt *time.Time `json:"updated_at"`
}

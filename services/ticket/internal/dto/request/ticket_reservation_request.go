package request

import (
	"github.com/google/uuid"
)

type TicketReservatioRequest struct {
	EventID  uuid.UUID `json:"event_id" validate:"required"`
	UserID   uuid.UUID `json:"user_id" validate:"required"`
	Quantity int       `json:"quantity" validate:"required"`
}

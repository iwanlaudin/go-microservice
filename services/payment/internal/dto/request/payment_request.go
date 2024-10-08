package request

import "github.com/google/uuid"

type PaymentRequest struct {
	ID uuid.UUID `json:"payment_id"`
}

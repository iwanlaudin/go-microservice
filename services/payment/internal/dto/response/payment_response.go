package response

import "time"

type PaymentResponse struct {
	ID        int        `json:"payment_id"`
	TicketID  int        `json:"ticket_id"`
	UserID    int        `json:"user_id"`
	Amount    float64    `json:"amount"`
	Status    string     `json:"payment_status"`
	Date      time.Time  `json:"payment_date"`
	CreatedBy *string    `json:"created_by,omitempty"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedBy *string    `json:"updated_by,omitempty"`
	UpdatedAt *time.Time `json:"updated_at"`
	IsDeleted bool       `json:"is_deleted"`
}

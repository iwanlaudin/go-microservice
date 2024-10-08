package models

import "time"

type Payment struct {
	ID        int        `db:"id" json:"payment_id"`
	TicketID  int        `db:"ticket_id" json:"ticket_id"`
	UserID    int        `db:"user_id" json:"user_id"`
	Amount    float64    `db:"amount" json:"amount"`
	Status    string     `db:"status" json:"payment_status"`
	Date      time.Time  `db:"date" json:"payment_date"`
	CreatedBy *string    `db:"created_by" json:"created_by,omitempty"`
	CreatedAt *time.Time `db:"created_at" json:"created_at"`
	UpdatedBy *string    `db:"updated_by" json:"updated_by,omitempty"`
	UpdatedAt *time.Time `db:"updated_at" json:"updated_at"`
	IsDeleted bool       `db:"is_deleted" json:"is_deleted"`
}

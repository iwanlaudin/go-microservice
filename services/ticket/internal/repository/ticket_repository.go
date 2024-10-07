package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/iwanlaudin/go-microservice/pkg/common/helpers"
	"github.com/iwanlaudin/go-microservice/services/ticket/internal/models"
	"github.com/jmoiron/sqlx"
)

type TicketRepository struct {
}

func NewTicketRepository() *TicketRepository {
	return &TicketRepository{}
}

func (repository *TicketRepository) AddTicket(ctx context.Context, db *sqlx.DB, ticket *models.Ticket) error {
	query := `
		INSERT INTO "tickets" (
            id, event_id, user_id, quantity, status, created_at, created_by
        ) VALUES (
			:id, :event_id, :user_id, :quantity, :status, :created_at, :created_by
        )
	`

	_, err := db.NamedExecContext(ctx, query, ticket)
	if err != nil {
		return helpers.CustomError("Failed to add ticket")
	}

	return nil
}

func (repository *TicketRepository) FindAllTicketByUserId(ctx context.Context, db *sqlx.DB, userId uuid.UUID) (*[]models.Ticket, error) {
	var tickets []models.Ticket

	query := `
		SELECT *
		FROM
			"tickets"
		WHERE
			is_deleted = false
		AND
			user_id = $1
	`
	err := db.SelectContext(ctx, &tickets, query, userId)
	if err != nil {
		if err == sql.ErrNoRows {
			return &tickets, nil
		}
		helpers.PanicIfError(err)
	}

	return &tickets, nil
}

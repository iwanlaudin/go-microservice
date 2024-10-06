package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/iwanlaudin/go-microservice/pkg/common/helpers"
	"github.com/iwanlaudin/go-microservice/services/event/internal/models"
	"github.com/jmoiron/sqlx"
)

type EventRepositoryImpl struct {
}

func NewEventRepository() EventRepository {
	return &EventRepositoryImpl{}
}

func (repository *EventRepositoryImpl) FindAll(ctx context.Context, db *sqlx.DB) (*[]models.Event, error) {
	var events []models.Event

	query := `
		SELECT * 
		FROM 
			"events"
		WHERE
			is_deleted = false
		ORDER BY
			name
		ASC
	`
	err := db.SelectContext(ctx, &events, query)
	if err != nil {
		if err == sql.ErrNoRows {
			return &events, nil
		}
		helpers.PanicIfError(err)
	}

	return &events, nil
}

func (repository *EventRepositoryImpl) FindById(ctx context.Context, db *sqlx.DB, id uuid.UUID) (*models.Event, error) {
	var event models.Event

	query := `
		SELECT *
		FROM
			"events"
		WHERE
			is_deleted = false AND id = $1
	`
	err := db.GetContext(ctx, &event, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return &event, nil
		}
		helpers.PanicIfError(err)
	}
	return &event, nil
}

func (repository *EventRepositoryImpl) AddEvent(ctx context.Context, db *sqlx.DB, event *models.Event) error {
	query := `
		INSERT INTO "events" (
			id, name, date, location, available_tickets, created_at, created_by
		) VALUES (
		 	:id, :name, :date, :location, :available_tickets, :created_at, :created_by
		)
	`
	_, err := db.NamedExecContext(ctx, query, event)
	if err != nil {
		helpers.PanicIfError(err)
	}
	return nil
}

func (repository *EventRepositoryImpl) AddEventRange(ctx context.Context, tx *sqlx.Tx, events *[]models.Event) error {
	if events == nil || len(*events) == 0 {
		return helpers.CustomError("No event to add")
	}

	query := `
        INSERT INTO "events" (
            id, name, date, location, available_tickets, created_at, created_by
        ) VALUES (
            :id, :title, :date, :location, :available_tickets, :created_at, :created_by
        )
    `
	for _, event := range *events {
		_, err := tx.NamedExecContext(ctx, query, event)
		if err != nil {
			return helpers.CustomError("Failed to add event")
		}
	}

	return nil
}

func (respository *EventRepositoryImpl) UpdateEvent(ctx context.Context, db *sqlx.DB, event *models.Event) error {
	if event == nil {
		return helpers.CustomError("No event to update")
	}

	query := `
		UPDATE "events"
		SET
			name = :name,
			date = :date,
			location = :location,
			available_tickets = :available_tickets
		WHERE
			is_deleted = false 
		AND
			id = :id
	`
	_, err := db.NamedExecContext(ctx, query, event)
	helpers.PanicIfError(err)

	return nil
}

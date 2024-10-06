package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/iwanlaudin/go-microservice/services/events/internal/models"
	"github.com/jmoiron/sqlx"
)

type EventRepository interface {
	FindAll(ctx context.Context, db *sqlx.DB) (*[]models.Event, error)
	FindById(ctx context.Context, db *sqlx.DB, id uuid.UUID) (*models.Event, error)
	AddEvent(ctx context.Context, db *sqlx.DB, event *models.Event) error
	AddEventRange(ctx context.Context, tx *sqlx.Tx, events *[]models.Event) error
	UpdateEvent(ctx context.Context, db *sqlx.DB, event *models.Event) error
}

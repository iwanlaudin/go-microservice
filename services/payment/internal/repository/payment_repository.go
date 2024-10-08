package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/iwanlaudin/go-microservice/pkg/common/api"
	"github.com/iwanlaudin/go-microservice/services/payment/internal/models"
	"github.com/jmoiron/sqlx"
)

type PaymentRepository struct {
}

func NewPaymentRepository() *PaymentRepository {
	return &PaymentRepository{}
}

func (repository *PaymentRepository) FindAll(ctx context.Context, db *sqlx.DB) (*[]models.Payment, error) {
	var payments []models.Payment

	query := `
		SELECT *
		FROM
			"payments"
		WHERE
			is_deleted = false
		AND
			user_id = $1
	`
	userId := api.UserIDFromContext(ctx)
	err := db.SelectContext(ctx, &payments, query, userId)
	if err != nil {
		return nil, err
	}
	return &payments, nil
}

func (repository *PaymentRepository) FindById(ctx context.Context, db *sqlx.DB, paymentId uuid.UUID) (*models.Payment, error) {
	var payment models.Payment

	query := `
		SELECT *
		FROM
			"payments"
		WHERE
			is_deleted = false
		AND
			id = $1
	`
	err := db.GetContext(ctx, payment, query, paymentId)
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

func (repository *PaymentRepository) Update(ctx context.Context, db *sqlx.DB, payment *models.Payment) error {
	query := `
		UPDATE "payments"
		SET
			status = :name,
			date = :date
		WHERE
			is_deleted = false 
		AND
			id = :id
	`
	_, err := db.NamedExecContext(ctx, query, payment)
	if err != nil {
		return err
	}
	return nil
}

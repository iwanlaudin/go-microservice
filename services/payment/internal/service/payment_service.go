package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/iwanlaudin/go-microservice/pkg/common/helpers"
	"github.com/iwanlaudin/go-microservice/services/payment/internal/dto/response"
	"github.com/iwanlaudin/go-microservice/services/payment/internal/repository"
	"github.com/jmoiron/sqlx"
)

type PaymentService struct {
	PaymentRepository *repository.PaymentRepository
	DB                *sqlx.DB
}

func NewPaymentService(paymentRepository *repository.PaymentRepository, db *sqlx.DB) *PaymentService {
	return &PaymentService{
		PaymentRepository: paymentRepository,
		DB:                db,
	}
}

func (service *PaymentService) FindAll(ctx context.Context) (*[]response.PaymentResponse, error) {
	var paymentResponse []response.PaymentResponse

	results, err := service.PaymentRepository.FindAll(ctx, service.DB)
	helpers.PanicIfError(err)

	for _, payment := range *results {
		paymentResponse = append(paymentResponse, response.PaymentResponse{
			ID:        payment.ID,
			TicketID:  payment.TicketID,
			UserID:    payment.UserID,
			Amount:    payment.Amount,
			Date:      payment.Date,
			CreatedBy: payment.CreatedBy,
			CreatedAt: payment.CreatedAt,
			UpdatedBy: payment.CreatedBy,
			UpdatedAt: payment.UpdatedAt,
		})
	}

	return &paymentResponse, nil
}

func (service *PaymentService) FindById(ctx context.Context, paymentId uuid.UUID) (*response.PaymentResponse, error) {
	result, err := service.PaymentRepository.FindById(ctx, service.DB, paymentId)
	helpers.PanicIfError(err)

	return &response.PaymentResponse{
		ID:        result.ID,
		TicketID:  result.TicketID,
		UserID:    result.UserID,
		Amount:    result.Amount,
		Date:      result.Date,
		CreatedBy: result.CreatedBy,
		CreatedAt: result.CreatedAt,
		UpdatedBy: result.CreatedBy,
		UpdatedAt: result.UpdatedAt,
	}, nil
}

func (service *PaymentService) Update(ctx context.Context, paymentId uuid.UUID) error {
	payment, err := service.PaymentRepository.FindById(ctx, service.DB, paymentId)
	helpers.PanicIfError(err)

	now := time.Now().UTC()
	payment.Date = now
	payment.Status = "paid"

	err = service.PaymentRepository.Update(ctx, service.DB, payment)
	helpers.PanicIfError(err)

	return nil
}

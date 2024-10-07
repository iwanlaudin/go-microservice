package consumer

import (
	"context"

	"github.com/iwanlaudin/go-microservice/pkg/common/logger"
	"github.com/iwanlaudin/go-microservice/pkg/rabbitmq"
	"github.com/jmoiron/sqlx"
)

type MessageConsumer struct {
	RabbitMQ *rabbitmq.RabbitMQ
	DB       *sqlx.DB
	Logger   logger.Logger
}

func NewMessageConsumer(rabbitMQ *rabbitmq.RabbitMQ, db *sqlx.DB, logger logger.Logger) *MessageConsumer {
	return &MessageConsumer{
		RabbitMQ: rabbitMQ,
		Logger:   logger,
		DB:       db,
	}
}

func (mc *MessageConsumer) StartConsuming(ctx context.Context) error {
	err := mc.RabbitMQ.ConsumeMessages(ctx, "ticket.reserved", mc.handleTicketReserved)
	if err != nil {
		mc.Logger.Error("Failed to consume ticket.reserved queue", logger.Error(err))
		return err
	}
	return nil
}

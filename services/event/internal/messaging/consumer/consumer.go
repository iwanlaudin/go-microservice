package consumer

import (
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

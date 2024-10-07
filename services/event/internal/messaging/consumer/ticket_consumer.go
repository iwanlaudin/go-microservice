package consumer

import (
	"context"
	"encoding/json"

	"github.com/iwanlaudin/go-microservice/pkg/common/logger"
	"github.com/iwanlaudin/go-microservice/services/event/internal/dto/messaging/consumer"
)

func (mc *MessageConsumer) StartConsuming(ctx context.Context) error {
	err := mc.RabbitMQ.ConsumeMessages(ctx, "ticket.reserved", mc.handleTicketReserved)
	return err
}

func (mc *MessageConsumer) handleTicketReserved(ctx context.Context, message []byte) {
	var reservationMsg consumer.TicketReservedMessage
	err := json.Unmarshal(message, &reservationMsg)
	if err != nil {
		mc.Logger.Error("Failed to unmarshal message", logger.Error(err))
		return
	}

	err = mc.processTicketReservation(ctx, reservationMsg)
	if err != nil {
		mc.Logger.Error("Failed to process ticket reservation", logger.Error(err))
		return
	}

	mc.Logger.Info("Successfully processed ticket reservation")
}

func (mc *MessageConsumer) processTicketReservation(ctx context.Context, msg consumer.TicketReservedMessage) error {
	mc.Logger.Info("Processing ticket reservation", logger.String("eventID", msg.EventID.String()))

	query := `
		UPDATE "events"
		SET
			available_tickets = available_tickets - :quantity
		WHERE
			is_deleted = false 
		AND
			id = :id
		AND
			available_tickets >= :quantity
	`
	params := map[string]interface{}{
		"id":       msg.EventID,
		"quantity": msg.Quantity,
	}
	_, err := mc.DB.NamedExecContext(ctx, query, params)
	if err != nil {
		return err
	}

	return nil
}

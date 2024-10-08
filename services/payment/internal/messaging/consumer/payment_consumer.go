package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/iwanlaudin/go-microservice/pkg/common/logger"
)

func (mc *MessageConsumer) handleTicketReserved(ctx context.Context, message []byte) {
	var reservationMsg map[string]interface{}
	err := json.Unmarshal(message, &reservationMsg)
	if err != nil {
		mc.Logger.Error("Failed to unmarshal message", logger.Error(err))
		return
	}

	reservationMsg["id"], _ = uuid.NewV7()
	reservationMsg["amount"] = 50.0
	reservationMsg["created_at"] = time.Now().UTC()
	reservationMsg["status"] = "unpaid"

	fmt.Printf("reservationMsg: %v\n", reservationMsg)

	err = mc.processTicketReservation(ctx, reservationMsg)
	if err != nil {
		mc.Logger.Error("Failed to process ticket reservation", logger.Error(err))
		return
	}

	mc.Logger.Info("Successfully processed ticket reservation")
}

func (mc *MessageConsumer) processTicketReservation(ctx context.Context, msg map[string]interface{}) error {
	mc.Logger.Info("Processing ticket reservation", logger.String("ticket_id", msg["ticket_id"].(string)))

	query := `
		INSERT INTO "payments" (
            id, ticket_id, user_id, amount, status, created_at
        ) VALUES (
            :id, :ticket_id, :user_id, :amount, :status, :created_at
        )
	`
	_, err := mc.DB.NamedExecContext(ctx, query, msg)
	if err != nil {
		return err
	}

	return nil
}

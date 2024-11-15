package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/iwanlaudin/go-microservice/pkg/common/api"
	"github.com/iwanlaudin/go-microservice/pkg/common/helpers"
	"github.com/iwanlaudin/go-microservice/pkg/common/logger"
	"github.com/iwanlaudin/go-microservice/pkg/rabbitmq"
	redisClient "github.com/iwanlaudin/go-microservice/pkg/redis"
	"github.com/iwanlaudin/go-microservice/services/ticket/internal/dto/messaging/publisher"
	"github.com/iwanlaudin/go-microservice/services/ticket/internal/dto/request"
	"github.com/iwanlaudin/go-microservice/services/ticket/internal/dto/response"
	"github.com/iwanlaudin/go-microservice/services/ticket/internal/models"
	"github.com/iwanlaudin/go-microservice/services/ticket/internal/repository"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

type TicketService struct {
	TicketRepository repository.TicketRepository
	DB               *sqlx.DB
	RedisClient      *redisClient.RedisClient
	HttpClient       *http.Client
	RabbitMQ         *rabbitmq.RabbitMQ
	Log              logger.Logger
	EventServiceURL  string
}

func NewTicketService(ticketRepository repository.TicketRepository, db *sqlx.DB, redisClient *redisClient.RedisClient, rabbitMQ *rabbitmq.RabbitMQ, log logger.Logger, eventServiceUrl string) *TicketService {
	return &TicketService{
		TicketRepository: ticketRepository,
		DB:               db,
		RedisClient:      redisClient,
		HttpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		RabbitMQ:        rabbitMQ,
		Log:             log,
		EventServiceURL: eventServiceUrl,
	}
}

func (service *TicketService) ReserveTicket(ctx context.Context, authToken string, request *request.TicketReservatioRequest) (*models.Ticket, error) {
	eventCache, err := service.getOrCreateEventCache(ctx, request.EventID, authToken)
	if err != nil {
		return nil, helpers.CustomError("Failed to get event info: %w", err)
	}

	if eventCache.Date.Before(time.Now()) {
		return nil, helpers.CustomError("Event has already passed")
	}

	if eventCache.AvailableTickets < request.Quantity {
		return nil, helpers.CustomError("Not enough tickets available")
	}

	eventCache.AvailableTickets -= request.Quantity
	err = service.createEventCache(ctx, eventCache)
	if err != nil {
		return nil, helpers.CustomError("Failed to update event cache: %w", err)
	}

	now := time.Now().UTC()
	eventId, err := uuid.NewV7()
	if err != nil {
		return nil, helpers.CustomError("Failed to generate event ID")
	}

	userIdStr := request.UserID.String()
	ticket := &models.Ticket{
		ID:        eventId,
		EventID:   request.EventID,
		UserID:    request.UserID,
		Quantity:  request.Quantity,
		Status:    "reserved",
		CreatedAt: &now,
		CreatedBy: &userIdStr,
	}

	err = service.TicketRepository.AddTicket(ctx, service.DB, ticket)
	if err != nil {
		eventCache.AvailableTickets += request.Quantity
		_ = service.createEventCache(ctx, eventCache)
		return nil, helpers.CustomError("Failed to create ticket: %w", err)
	}

	_ = service.publishTicketReservedCreated(ctx, ticket)

	return ticket, nil
}

func (service *TicketService) GetAllTicketByUser(ctx context.Context) (*[]response.TicketResponse, error) {
	var ticketResponse []response.TicketResponse

	userId := api.UserIDFromContext(ctx)

	tickets, err := service.TicketRepository.FindAllTicketByUserId(ctx, service.DB, userId)
	if err != nil {
		return nil, helpers.CustomError(err.Error())
	}

	for _, ticket := range *tickets {
		ticketResponse = append(ticketResponse, response.TicketResponse{
			ID:        ticket.ID,
			EventID:   ticket.EventID,
			UserID:    ticket.UserID,
			Quantity:  ticket.Quantity,
			Status:    ticket.Status,
			CreatedBy: ticket.CreatedBy,
			CreatedAt: ticket.CreatedAt,
		})
	}

	return &ticketResponse, nil
}

func (service *TicketService) publishTicketReservedCreated(ctx context.Context, ticket *models.Ticket) error {
	reservationMsg := publisher.TicketReservedEvent{
		ID:       ticket.ID,
		EventID:  ticket.EventID,
		UserID:   ticket.UserID,
		Quantity: ticket.Quantity,
	}

	reservationMsgJSON, err := json.Marshal(reservationMsg)
	if err != nil {
		service.Log.Error("Failed to encode json", logger.Error(err))
	}

	err = service.publishWithRetry(ctx, "ticket.reserved", "ticket.created", reservationMsgJSON)
	if err == nil {
		service.Log.Info("successfully to publish ticket created event after retries")
	}

	err = service.publishWithRetry(ctx, "ticket.payment", "ticket.payment", reservationMsgJSON)
	if err == nil {
		service.Log.Info("successfully to publish ticket created payment after retries")
	}

	if err != nil {
		return err
	}

	return nil
}

func (service *TicketService) publishWithRetry(ctx context.Context, exchange, routingKey string, body []byte) error {
	maxRetries := 3
	for i := 0; i < maxRetries; i++ {
		err := service.RabbitMQ.PublishMessage(ctx, exchange, routingKey, body)
		if err == nil {
			return nil
		}
		service.Log.Warn("Failed to publish message, retrying", logger.Error(err))
		time.Sleep(time.Second * time.Duration(i+1))
	}
	return helpers.CustomError("Failed to publish message after %d attempts", maxRetries)
}

func (service *TicketService) getOrCreateEventCache(ctx context.Context, eventId uuid.UUID, authToken string) (*models.Event, error) {
	eventCache, err := service.RedisClient.Get(ctx, eventId.String())
	if err != nil {
		if !errors.Is(err, redis.Nil) {
			return nil, helpers.CustomError("Event not found in cache: %w", err)
		}
	}

	var eventResult *models.Event

	if err := json.Unmarshal([]byte(eventCache), &eventResult); err == nil {
		return eventResult, nil
	}

	eventResult, err = service.fetchEventFromService(ctx, eventId, authToken)
	if err != nil {
		return nil, helpers.CustomError("Failed to fetch event from service: %w", err)
	}

	err = service.RedisClient.Set(ctx, eventId.String(), eventResult, time.Hour*1)
	if err != nil {
		return nil, helpers.CustomError("Failed to set event cache: %v\n", err)
	}

	return eventResult, nil
}

func (service *TicketService) createEventCache(ctx context.Context, event *models.Event) error {
	err := service.RedisClient.Set(ctx, event.ID.String(), event, time.Hour*1)
	if err != nil {
		return helpers.CustomError("Failed to set event cache: %v\n", err)
	}

	return nil
}

func (service *TicketService) fetchEventFromService(ctx context.Context, eventID uuid.UUID, authToken string) (*models.Event, error) {
	url := fmt.Sprintf("%s/api/events/%s", service.EventServiceURL, eventID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, helpers.CustomError("Failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+authToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := service.HttpClient.Do(req)
	if err != nil {
		return nil, helpers.CustomError("Failed to send request: %w", err)
	}
	defer resp.Body.Close()

	var eventResponse response.EventResponse

	if resp.StatusCode != http.StatusOK {
		return nil, helpers.CustomError("Unexpected status code: %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&eventResponse); err != nil {
		return nil, helpers.CustomError("Failed to decode response: %w", err)
	}

	if eventResponse.Items == nil {
		return nil, helpers.CustomError("Event items is nil or empty")
	}

	eventResult := models.Event{
		ID:               eventResponse.Items.ID,
		Name:             eventResponse.Items.Name,
		Date:             eventResponse.Items.Date,
		AvailableTickets: eventResponse.Items.AvailableTickets,
	}

	return &eventResult, nil
}

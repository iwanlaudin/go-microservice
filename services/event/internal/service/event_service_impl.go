package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/iwanlaudin/go-microservice/pkg/common/api"
	"github.com/iwanlaudin/go-microservice/pkg/common/helpers"
	"github.com/iwanlaudin/go-microservice/services/event/internal/dto/request"
	"github.com/iwanlaudin/go-microservice/services/event/internal/dto/response"
	"github.com/iwanlaudin/go-microservice/services/event/internal/models"
	"github.com/iwanlaudin/go-microservice/services/event/internal/repository"
	"github.com/jmoiron/sqlx"
)

type EventServiceImpl struct {
	EventRepository repository.EventRepository
	DB              *sqlx.DB
}

func NewEventService(eventRepository repository.EventRepository, db *sqlx.DB) EventService {
	return &EventServiceImpl{
		EventRepository: eventRepository,
		DB:              db,
	}
}

func (service *EventServiceImpl) CreateEvent(ctx context.Context, request *request.CreateEventRequest) (*response.EventResponse, error) {
	userContext := api.UserFromContext(ctx)
	eventId, err := uuid.NewV7()
	if err != nil {
		return nil, helpers.CustomError("Failed to generate event ID")
	}

	now := time.Now().UTC()
	event := models.Event{
		ID:               eventId,
		Name:             request.Name,
		Date:             request.Date,
		Location:         request.Location,
		AvailableTickets: request.AvailableTickets,
		CreatedBy:        &userContext.ID,
		CreatedAt:        &now,
	}

	if err := service.EventRepository.AddEvent(ctx, service.DB, &event); err != nil {
		return nil, helpers.CustomError(err.Error())
	}

	return &response.EventResponse{
		ID:               event.ID,
		Name:             event.Name,
		Date:             event.Date,
		Location:         event.Location,
		AvailableTickets: event.AvailableTickets,
		CreatedBy:        event.CreatedBy,
		CreatedAt:        event.CreatedAt,
	}, nil
}

func (service *EventServiceImpl) UpdateEvent(ctx context.Context, request *request.UpdateEventRequest) (*response.EventResponse, error) {
	userContext := api.UserFromContext(ctx)

	event, err := service.EventRepository.FindById(ctx, service.DB, request.ID)
	if err != nil {
		return nil, helpers.CustomError(err.Error())
	}

	now := time.Now().UTC()
	event.Name = request.Name
	event.Date = request.Date
	event.Location = request.Location
	event.AvailableTickets = request.AvailableTickets
	event.UpdatedBy = &userContext.ID
	event.UpdatedAt = &now

	if err := service.EventRepository.UpdateEvent(ctx, service.DB, event); err != nil {
		return nil, helpers.CustomError(err.Error())
	}

	return &response.EventResponse{
		ID:               event.ID,
		Name:             event.Name,
		Date:             event.Date,
		Location:         event.Location,
		AvailableTickets: event.AvailableTickets,
		CreatedBy:        event.CreatedBy,
		CreatedAt:        event.CreatedAt,
		UpdatedBy:        event.UpdatedBy,
		UpdatedAt:        event.UpdatedAt,
	}, nil
}

func (service *EventServiceImpl) GetAllEvent(ctx context.Context) (*[]response.EventResponse, error) {
	events, err := service.EventRepository.FindAll(ctx, service.DB)
	if err != nil {
		return nil, helpers.CustomError(err.Error())
	}

	var eventResponses []response.EventResponse

	for _, event := range *events {
		eventResponses = append(eventResponses, response.EventResponse{
			ID:               event.ID,
			Name:             event.Name,
			Date:             event.Date,
			Location:         event.Location,
			AvailableTickets: event.AvailableTickets,
			CreatedBy:        event.CreatedBy,
			CreatedAt:        event.CreatedAt,
			UpdatedBy:        event.UpdatedBy,
			UpdatedAt:        event.UpdatedAt,
		})
	}

	return &eventResponses, nil
}

func (service *EventServiceImpl) GetEventById(ctx context.Context, request *request.GetEventRequest) (*response.EventResponse, error) {
	eventId, err := helpers.ConvertStringToUUID(request.ID)
	helpers.PanicIfError(err)

	event, err := service.EventRepository.FindById(ctx, service.DB, eventId)
	if err != nil {
		return nil, helpers.CustomError(err.Error())
	}

	return &response.EventResponse{
		ID:               event.ID,
		Name:             event.Name,
		Date:             event.Date,
		Location:         event.Location,
		AvailableTickets: event.AvailableTickets,
		CreatedBy:        event.CreatedBy,
		CreatedAt:        event.CreatedAt,
		UpdatedBy:        event.UpdatedBy,
		UpdatedAt:        event.UpdatedAt,
	}, nil
}

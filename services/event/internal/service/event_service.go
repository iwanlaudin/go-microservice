package service

import (
	"context"

	"github.com/iwanlaudin/go-microservice/services/event/internal/dto/request"
	"github.com/iwanlaudin/go-microservice/services/event/internal/dto/response"
)

type EventService interface {
	CreateEvent(ctx context.Context, request *request.CreateEventRequest) (*response.EventResponse, error)
	UpdateEvent(ctx context.Context, request *request.UpdateEventRequest) (*response.EventResponse, error)
	GetAllEvent(ctx context.Context) (*[]response.EventResponse, error)
	GetEventById(ctx context.Context, request *request.GetEventRequest) (*response.EventResponse, error)
}

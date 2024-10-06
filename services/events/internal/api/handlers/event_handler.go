package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/iwanlaudin/go-microservice/pkg/common/api"
	"github.com/iwanlaudin/go-microservice/pkg/common/helpers"
	"github.com/iwanlaudin/go-microservice/services/events/internal/dto/request"
	"github.com/iwanlaudin/go-microservice/services/events/internal/service"
)

type EventHandler struct {
	EventService service.EventService
	Validate     *validator.Validate
}

func NewEventHandler(eventService service.EventService, validate *validator.Validate) *EventHandler {
	return &EventHandler{
		EventService: eventService,
		Validate:     validate,
	}
}

func (handler *EventHandler) GetAllEvent(writer http.ResponseWriter, req *http.Request) {
	eventResponse, err := handler.EventService.GetAllEvent(req.Context())
	if err != nil {
		api.NewAppResponse(err.Error(), http.StatusBadRequest).Err(writer)
		return
	}

	api.NewAppResponse("Successfully to find event", http.StatusOK).Ok(writer, eventResponse)
}

func (handler *EventHandler) CreateEvent(writer http.ResponseWriter, req *http.Request) {
	createEventRequest := request.CreateEventRequest{}
	helpers.ReadFromRequestBody(req, &createEventRequest)

	if err := handler.Validate.Struct(createEventRequest); err != nil {
		api.NewAppResponse("Invalid parameter", http.StatusBadRequest).ValidationErr(writer, err)
		return
	}

	eventResponse, err := handler.EventService.CreateEvent(req.Context(), &createEventRequest)
	if err != nil {
		api.NewAppResponse(err.Error(), http.StatusBadRequest).Err(writer)
		return
	}
	api.NewAppResponse("Successfully to created event", http.StatusCreated).Ok(writer, eventResponse)
}

func (handler *EventHandler) GetEventById(writer http.ResponseWriter, req *http.Request) {
	getEventRequest := request.GetEventRequest{ID: chi.URLParam(req, "eventId")}

	eventResponse, err := handler.EventService.GetEventById(req.Context(), &getEventRequest)
	if err != nil {
		api.NewAppResponse(err.Error(), http.StatusBadRequest).Err(writer)
		return
	}
	api.NewAppResponse("Successfully to find event", http.StatusOK).Ok(writer, eventResponse)
}

func (handler *EventHandler) UpdateEvent(writer http.ResponseWriter, req *http.Request) {
	updateEventRequest := request.UpdateEventRequest{}
	helpers.ReadFromRequestBody(req, &updateEventRequest)

	eventId := chi.URLParam(req, "eventId")
	id, err := helpers.ConvertStringToUUID(eventId)
	helpers.PanicIfError(err)

	updateEventRequest.ID = id

	if err := handler.Validate.Struct(updateEventRequest); err != nil {
		api.NewAppResponse("Invalid parameter", http.StatusBadRequest).ValidationErr(writer, err)
		return
	}

	eventResponse, err := handler.EventService.UpdateEvent(req.Context(), &updateEventRequest)
	if err != nil {
		api.NewAppResponse(err.Error(), http.StatusBadRequest).Err(writer)
		return
	}
	api.NewAppResponse("Successfully to update event", http.StatusOK).Ok(writer, eventResponse)
}

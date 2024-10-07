package handlers

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/iwanlaudin/go-microservice/pkg/common/api"
	"github.com/iwanlaudin/go-microservice/pkg/common/helpers"
	"github.com/iwanlaudin/go-microservice/services/ticket/internal/dto/request"
	"github.com/iwanlaudin/go-microservice/services/ticket/internal/service"
)

type TicketHandler struct {
	TicketService service.TicketService
	Validate      *validator.Validate
}

func NewTicketHandler(ticketService service.TicketService, validate *validator.Validate) *TicketHandler {
	return &TicketHandler{
		TicketService: ticketService,
		Validate:      validate,
	}
}

func (handler *TicketHandler) ReserveTicket(writer http.ResponseWriter, req *http.Request) {
	reservationTicketRequest := request.TicketReservatioRequest{}
	helpers.ReadFromRequestBody(req, &reservationTicketRequest)

	eventIdReq := chi.URLParam(req, "eventId")
	eventId, _ := helpers.ConvertStringToUUID(eventIdReq)

	userId := api.UserIDFromContext(req.Context())

	reservationTicketRequest.UserID = userId
	reservationTicketRequest.EventID = eventId

	if err := handler.Validate.Struct(reservationTicketRequest); err != nil {
		api.NewAppResponse("Invalid parameter", http.StatusBadRequest).ValidationErr(writer, err)
		return
	}

	authToken := strings.Fields(req.Header.Get("Authorization"))
	ticketResponse, err := handler.TicketService.ReserveTicket(req.Context(), authToken[1], &reservationTicketRequest)
	if err != nil {
		api.NewAppResponse(err.Error(), http.StatusBadRequest).Err(writer)
		return
	}
	api.NewAppResponse("Successfully to update event", http.StatusOK).Ok(writer, ticketResponse)
}

func (handler *TicketHandler) GetAllTicketByUser(writer http.ResponseWriter, req *http.Request) {
	ticketResponse, err := handler.TicketService.GetAllTicketByUser(req.Context())
	if err != nil {
		api.NewAppResponse(err.Error(), http.StatusBadRequest).Err(writer)
		return
	}
	api.NewAppResponse("Successfully to update event", http.StatusOK).Ok(writer, ticketResponse)
}

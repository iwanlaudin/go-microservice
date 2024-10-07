package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/iwanlaudin/go-microservice/pkg/common/api"
	"github.com/iwanlaudin/go-microservice/services/ticket/internal/api/handlers"
)

func OrderRoute(handler *handlers.TicketHandler) chi.Router {
	r := chi.NewRouter()

	r.Route("/tickets", func(r chi.Router) {
		r.Use(api.AuthMiddleware)

		r.Get("/", handler.GetAllTicketByUser)
		r.Route("/{eventId}", func(r chi.Router) {
			r.Post("/", handler.ReserveTicket)
			r.Get("/", func(w http.ResponseWriter, r *http.Request) {

			})
		})
	})
	return r
}

package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/iwanlaudin/go-microservice/pkg/common/api"
	"github.com/iwanlaudin/go-microservice/services/events/internal/api/handlers"
)

func EventRouter(handler *handlers.EventHandler) chi.Router {
	r := chi.NewRouter()

	r.Route("/events", func(r chi.Router) {
		r.Use(api.AuthMiddleware)

		r.Post("/", handler.CreateEvent)
		r.Get("/", handler.GetAllEvent)

		r.Route("/{eventId}", func(r chi.Router) {
			r.Get("/", handler.GetEventById)
			r.Put("/", handler.UpdateEvent)
		})
	})

	return r
}

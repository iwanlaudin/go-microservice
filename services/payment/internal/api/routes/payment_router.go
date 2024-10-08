package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/iwanlaudin/go-microservice/pkg/common/api"
	"github.com/iwanlaudin/go-microservice/services/payment/internal/api/handlers"
)

func PaymentRouter(handler *handlers.HandlerPayment) chi.Router {
	r := chi.NewRouter()

	r.Route("/payments", func(r chi.Router) {
		r.Use(api.AuthMiddleware)

		r.Get("/", handler.FindAll)
		r.Route("/{paymentId}", func(r chi.Router) {
			r.Get("/", handler.FindById)
			r.Put("/", handler.Update)
		})
	})

	return r
}

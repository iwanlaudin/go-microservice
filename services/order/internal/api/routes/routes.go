package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/iwanlaudin/go-microservice/pkg/common/api"
	"github.com/iwanlaudin/go-microservice/pkg/common/auth"
)

func SetupRoutes(r chi.Router) {
	r.Use(api.ErrorLogger)

	r.Route("/orders", func(r chi.Router) {
		r.Use(auth.AuthMiddleware)

		// ... other routes
	})
}

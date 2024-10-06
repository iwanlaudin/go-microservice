package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/iwanlaudin/go-microservice/pkg/common/api"
	"github.com/iwanlaudin/go-microservice/services/authentication/internal/api/handlers"
)

func AuthRouter(handler *handlers.AuthHandler) chi.Router {
	r := chi.NewRouter()

	r.Post("/auth/sign-in", handler.SignIn)
	r.Post("/auth/sign-up", handler.SignUp)
	r.Route("/auth", func(r chi.Router) {
		r.Use(api.AuthMiddleware)
		r.Post("/refresh-token", handler.RefreshToken)
		r.Get("/me", handler.GetMe)
	})

	return r
}

package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/iwanlaudin/go-microservice/pkg/common/api"
	"github.com/iwanlaudin/go-microservice/services/authentication/internal/api/handlers"
)

func AuthRouter(h *handlers.AuthHandler) chi.Router {
	r := chi.NewRouter()

	r.Post("/sign-in", h.SignIn)
	r.Post("/sign-up", h.CreateUser)
	r.With(api.AuthMiddleware).Post("/refresh-token", h.RefreshToken)
	r.With(api.AuthMiddleware).Get("/me", h.Me)

	return r
}

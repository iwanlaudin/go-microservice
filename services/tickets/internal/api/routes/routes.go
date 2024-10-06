package routes

import (
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/iwanlaudin/go-microservice/pkg/common/api"
	"github.com/jmoiron/sqlx"
)

func NewRoute(db *sqlx.DB, validate *validator.Validate) *chi.Mux {
	// Initialize Router
	r := chi.NewRouter()

	// Base Middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Cutome Middleware
	r.Use(api.TimeoutMiddleware(15 * time.Second))
	r.Use(api.ErrorLogger)

	// Initialize Repository, Service and Handler

	// Apps router
	r.Mount("/api", OrderRoute())

	return r
}

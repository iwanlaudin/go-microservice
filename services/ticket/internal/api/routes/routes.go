package routes

import (
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/iwanlaudin/go-microservice/pkg/common/api"
	"github.com/iwanlaudin/go-microservice/pkg/common/config"
	"github.com/iwanlaudin/go-microservice/pkg/common/logger"
	"github.com/iwanlaudin/go-microservice/pkg/rabbitmq"
	"github.com/iwanlaudin/go-microservice/pkg/redis"
	"github.com/iwanlaudin/go-microservice/services/ticket/internal/api/handlers"
	"github.com/iwanlaudin/go-microservice/services/ticket/internal/repository"
	"github.com/iwanlaudin/go-microservice/services/ticket/internal/service"
	"github.com/jmoiron/sqlx"
)

func NewRoute(db *sqlx.DB, redis *redis.RedisClient, rabbitMQ *rabbitmq.RabbitMQ, validate *validator.Validate, log logger.Logger, cfg *config.Config) *chi.Mux {
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
	ticketRepository := repository.NewTicketRepository()
	ticketService := service.NewTicketService(*ticketRepository, db, redis, rabbitMQ, log, cfg.EventServiceURL)
	ticketHandler := handlers.NewTicketHandler(*ticketService, validate)

	// Apps router
	r.Mount("/api", OrderRoute(ticketHandler))

	return r
}

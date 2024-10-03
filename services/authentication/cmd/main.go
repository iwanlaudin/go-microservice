package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/iwanlaudin/go-microservice/pkg/common/api"
	"github.com/iwanlaudin/go-microservice/pkg/common/config"
	"github.com/iwanlaudin/go-microservice/pkg/common/database"
	"github.com/iwanlaudin/go-microservice/pkg/common/logger"
	"github.com/iwanlaudin/go-microservice/services/authentication/internal/api/routes"
)

func main() {
	// Load .env Configuration
	cfg := config.Load()

	// Initialize logger
	log := logger.New(cfg.LogLevel)

	// Initialize Database Connection
	db, err := database.NewConnection(cfg.AuthDbURL)
	if err != nil {
		log.Fatal("Failed to connect to database", logger.Error(err))
	}
	defer db.Close()

	// Run Database Migrations
	if err := database.RunMigrations(db); err != nil {
		log.Fatal("Failed to run database migrations", logger.Error(err))
	}

	// Initialize Router
	r := chi.NewRouter()

	// Base Middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Cutome Middleware
	// Set global timeout to 15 seconds
	r.Use(api.TimeoutMiddleware(15 * time.Second))

	// Route
	routes.SetupRoutes(r)

	// Configuration Server
	srv := &http.Server{
		Addr:         cfg.AuthServicePort,
		Handler:      r,
		ReadTimeout:  20 * time.Second,
		WriteTimeout: 20 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Run server
	go func() {
		log.Info("Starting Order service", logger.String("port", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Could not listen on", logger.String("port", srv.Addr), logger.Error(err))
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", logger.Error(err))
	}

	log.Info("Server exiting")
}

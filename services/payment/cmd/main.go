package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/iwanlaudin/go-microservice/pkg/common/config"
	"github.com/iwanlaudin/go-microservice/pkg/common/database"
	"github.com/iwanlaudin/go-microservice/pkg/common/logger"
	"github.com/iwanlaudin/go-microservice/services/payment/internal/api/routes"
)

func main() {
	// Load konfigurasi
	cfg := config.Load()

	// Inisialisasi logger
	log := logger.New(cfg.LogLevel)

	// Inisialisasi koneksi database
	db, err := database.NewConnection(cfg.PaymentDbURL)
	if err != nil {
		log.Fatal("Failed to connect to database", logger.Error(err))
	}
	defer db.Close()

	// Jalankan migrasi database
	if err := database.RunMigrations(db); err != nil {
		log.Fatal("Failed to run database migrations", logger.Error(err))
	}

	// Initialize validator
	validate := validator.New()

	// Inisialisasi router
	handler := routes.NewRouter(db, validate)

	// Konfigurasi server
	srv := &http.Server{
		Addr:         cfg.PaymentServicePort,
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	// Run server
	go func() {
		log.Info("Starting Payment service", logger.String("port", srv.Addr))
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

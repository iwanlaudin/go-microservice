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
	"github.com/iwanlaudin/go-microservice/pkg/rabbitmq"
	"github.com/iwanlaudin/go-microservice/pkg/redis"
	"github.com/iwanlaudin/go-microservice/services/event/internal/api/routes"
	"github.com/iwanlaudin/go-microservice/services/event/internal/messaging/consumer"
)

func main() {
	// Load .env Configuration
	cfg := config.Load()

	// Initialize logger
	log := logger.New(cfg.LogLevel)

	// Initialize Database Connection
	db, err := database.NewConnection(cfg.EventDbURL)
	if err != nil {
		log.Fatal("Failed to connect to database", logger.Error(err))
	}
	defer db.Close()

	// Initialize redis client
	redisClient, err := redis.NewRedisClient(cfg.RedisURL)
	if err != nil {
		log.Fatal("Failed to connect to redis", logger.Error(err))
	}
	defer redisClient.Client.Close()

	// Initialize RabbitMQ
	rabbitMQ, err := rabbitmq.NewRabbitMQ(cfg.RabbitMQURL)
	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ:", logger.Error(err))
	}
	defer rabbitMQ.Close()

	// Initialize validator
	validate := validator.New()

	// Initialize Router
	handler := routes.NewRouter(db, validate)

	// Configuration Server
	srv := &http.Server{
		Addr:         cfg.EventServicePort,
		Handler:      handler,
		ReadTimeout:  20 * time.Second,
		WriteTimeout: 20 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Initialize and start RabbitMQ consumer
	newConsumer := consumer.NewMessageConsumer(rabbitMQ, db, log)
	startRabbitMQConsumer(context.Background(), newConsumer)

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

func startRabbitMQConsumer(ctx context.Context, consumer *consumer.MessageConsumer) {
	go func() {
		err := consumer.StartConsuming(ctx)
		if err != nil {
			consumer.Logger.Error("Failed to consume RabbitMQ messages", logger.Error(err))
		}
	}()
}

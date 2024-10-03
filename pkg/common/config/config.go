package config

import (
	"errors"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

var (
	secretKey      []byte
	ErrNoSecretKey = errors.New("no secret key provided")
)

type Config struct {
	AuthDbURL          string
	OrderDbURL         string
	PaymentDbURL       string
	OrderServicePort   string
	PaymentServicePort string
	LogLevel           string
	SecretKey          string
	RabbitMQURL        string
}

func Load() *Config {
	root, _ := os.Getwd()
	envFile := filepath.Join(root, ".env")
	if err := godotenv.Load(envFile); err != nil {
		log.Printf("No .env file found, using environment variables")
	}

	config := &Config{
		AuthDbURL:          getEnv("ORDER_DB_URL", "postgres://user:password@authDb:5432/orderdb?sslmode=disable"),
		OrderDbURL:         getEnv("ORDER_DB_URL", "postgres://user:password@orderdb:5432/orderdb?sslmode=disable"),
		PaymentDbURL:       getEnv("PAYMENT_DB_URL", "postgres://user:password@orderdb:5432/paymentDb?sslmode=disable"),
		OrderServicePort:   getEnv("ORDER_SERVICE_PORT", ":8080"),
		PaymentServicePort: getEnv("PAYMENT_SERVICE_PORT", ":8081"),
		LogLevel:           getEnv("LOG_LEVEL", "info"),
		SecretKey:          getEnv("SECRET_KEY", "x-secret-key"),
		RabbitMQURL:        getEnv("RABBIT_MQ_URL", "amqp://guest:guest@localhost:5672/"),
	}

	secretKey = []byte(config.SecretKey)

	return config
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func GetSecretKey() ([]byte, error) {
	if len(secretKey) == 0 {
		return nil, ErrNoSecretKey
	}
	return secretKey, nil
}

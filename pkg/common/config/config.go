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
	EventDbURL         string
	TicketDbURL        string
	PaymentDbURL       string
	AuthServicePort    string
	EventServicePort   string
	TicketServicePort  string
	PaymentServicePort string
	LogLevel           string
	SecretKey          string
	RabbitMQURL        string
	RedisURL           string
	SmtpHost           string
	SmtpPort           string
	SmtpUsername       string
	SmtpPassword       string
}

func Load() *Config {
	root, _ := os.Getwd()
	envFile := filepath.Join(root, ".env")
	if err := godotenv.Load(envFile); err != nil {
		log.Printf("No .env file found, using environment variables")
	}

	config := &Config{
		AuthDbURL:          getEnv("AUTH_DB_URL", "postgres://user:password@authDb:5432/authDb?sslmode=disable"),
		EventDbURL:         getEnv("EVENT_DB_URL", "postgres://user:password@authDb:5432/eventDb?sslmode=disable"),
		TicketDbURL:        getEnv("TICKET_DB_URL", "postgres://user:password@orderdb:5432/orderdb?sslmode=disable"),
		PaymentDbURL:       getEnv("PAYMENT_DB_URL", "postgres://user:password@orderdb:5432/paymentDb?sslmode=disable"),
		AuthServicePort:    getEnv("AUTH_SERVICE_PORT", ":8080"),
		EventServicePort:   getEnv("EVENT_SERVICE_PORT", ":8081"),
		TicketServicePort:  getEnv("TICKET_SERVICE_PORT", ":8082"),
		PaymentServicePort: getEnv("PAYMENT_SERVICE_PORT", ":8083"),
		LogLevel:           getEnv("LOG_LEVEL", "info"),
		SecretKey:          getEnv("SECRET_KEY", "x-secret-key"),
		RabbitMQURL:        getEnv("RABBIT_MQ_URL", "amqp://guest:guest@localhost:5672/"),
		RedisURL:           getEnv("REDIS_URL", "redis://:<password>@localhost:6379/0"),
		SmtpHost:           getEnv("SMTP_HOST", ""),
		SmtpPort:           getEnv("SMTP_PORT", ""),
		SmtpUsername:       getEnv("SMTP_USERNAME", ""),
		SmtpPassword:       getEnv("SMTP_PASSWORD", ""),
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

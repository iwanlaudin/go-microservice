package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost             string
	DBPort             string
	DBUser             string
	DBPassword         string
	DBName             string
	DBSSLMode          string
	OrderServicePort   string
	PaymentServicePort string
	LogLevel           string
}

func Load() *Config {
	// Cari file .env
	envFile := filepath.Join(".env")
	if err := godotenv.Load(envFile); err != nil {
		log.Printf("No .env file found, using environment variables")
	}

	return &Config{
		DBHost:             getEnv("DB_HOST", "localhost"),
		DBPort:             getEnv("DB_PORT", "5432"),
		DBUser:             getEnv("DB_USER", "postgres"),
		DBPassword:         getEnv("DB_PASSWORD", ""),
		DBName:             getEnv("DB_NAME", "myproject"),
		DBSSLMode:          getEnv("DB_SSLMODE", "disable"),
		OrderServicePort:   getEnv("ORDER_SERVICE_PORT", "8080"),
		PaymentServicePort: getEnv("PAYMENT_SERVICE_PORT", "8081"),
		LogLevel:           getEnv("LOG_LEVEL", "info"),
	}
}

// getEnv mengambil nilai environment variable atau mengembalikan nilai default
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

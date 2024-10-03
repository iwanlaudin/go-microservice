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
	DBHost             string
	DBPort             string
	DBUser             string
	DBPassword         string
	DBName             string
	DBSSLMode          string
	OrderServicePort   string
	PaymentServicePort string
	LogLevel           string
	SecretKey          string
}

func Load() *Config {
	root, _ := os.Getwd()
	envFile := filepath.Join(root, ".env")
	if err := godotenv.Load(envFile); err != nil {
		log.Printf("No .env file found, using environment variables")
	}

	config := &Config{
		DBHost:             getEnv("DB_HOST", "localhost"),
		DBPort:             getEnv("DB_PORT", "5432"),
		DBUser:             getEnv("DB_USER", "postgres"),
		DBPassword:         getEnv("DB_PASSWORD", ""),
		DBName:             getEnv("DB_NAME", "myproject"),
		DBSSLMode:          getEnv("DB_SSLMODE", "disable"),
		OrderServicePort:   getEnv("ORDER_SERVICE_PORT", ":8080"),
		PaymentServicePort: getEnv("PAYMENT_SERVICE_PORT", ":8081"),
		LogLevel:           getEnv("LOG_LEVEL", "info"),
		SecretKey:          getEnv("SECRET_KEY", "x-secret-key"),
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

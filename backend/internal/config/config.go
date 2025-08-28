package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Port               string
	DatabaseURL        string
	JWTSecret          string
	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectURL  string
	OpenAIAPIKey       string
	OpenAIEndpoint     string
	LogLevel           logrus.Level
	Environment        string
}

func Load() (*Config, error) {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		// .env file is optional, continue without it
	}

	config := &Config{
		Port:               getEnv("PORT", "8080"),
		DatabaseURL:        getEnv("DATABASE_URL", "postgres://postgres:password@localhost:5432/factcheck?sslmode=disable"),
		JWTSecret:          getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
		GoogleClientID:     getEnv("GOOGLE_CLIENT_ID", ""),
		GoogleClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
		GoogleRedirectURL:  getEnv("GOOGLE_REDIRECT_URL", "http://localhost:3000/auth/callback"),
		OpenAIAPIKey:       getEnv("OPENAI_API_KEY", ""),
		OpenAIEndpoint:     getEnv("OPENAI_ENDPOINT", "https://api.openai.com/v1"),
		Environment:        getEnv("ENVIRONMENT", "development"),
	}

	// Parse log level
	logLevel, err := logrus.ParseLevel(getEnv("LOG_LEVEL", "info"))
	if err != nil {
		logLevel = logrus.InfoLevel
	}
	config.LogLevel = logLevel

	return config, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

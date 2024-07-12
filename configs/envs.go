package configs

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

type Config struct {
	DbUser          string
	DbPassword      string
	DbName          string
	DbHost          string
	DbPort          string
	RateLimit       int
	RateLimitWindow int
	Port            string
	JwtSecret       string
}

func NewEnv() *Config {

	// load .env file
	err := godotenv.Load()

	if err != nil {
		zap.L().Error("Error loading .env file")
	}

	zap.L().Info("Loaded .env file")
	return &Config{
		DbUser:          getEnv("POSTGRES_USER", "postgres"),
		DbPassword:      getEnv("POSTGRES_PASSWORD", "postgres"),
		DbName:          getEnv("POSTGRES_DB", "postgres"),
		DbHost:          getEnv("POSTGRES_HOST", "localhost"),
		DbPort:          getEnv("POSTGRES_PORT", "5432"),
		RateLimit:       getEnvInt("RATE_LIMITER_MAX_REQUESTS", 10),
		RateLimitWindow: getEnvInt("RATE_LIMITER_WINDOW", 10),
		Port:            getEnv("PORT", "3000"),
		JwtSecret:       getEnv("JWT_SECRET", "secret"),
	}
}

func getEnv(key string, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}

func getEnvInt(key string, defaultValue int) int {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}

	val, _ := strconv.Atoi(value)

	return val
}

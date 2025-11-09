package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config captures runtime configuration for the search service.
type Config struct {
	Env               string
	AppPort           string
	GroupServiceURL   string
	EnablePlayground  bool
	DefaultResultSize int64
	MaxResultSize     int64
	HTTPClientTimeout time.Duration
}

// Load reads configuration from environment variables (optionally via .env).
func Load() (*Config, error) {
	if os.Getenv("APP_ENV") == "" {
		_ = godotenv.Load(".env")
	}

	enablePlayground := getEnvAsBool("ENABLE_GRAPHQL_PLAYGROUND", true)
	defaultLimit := getEnvAsInt64("DEFAULT_SEARCH_LIMIT", 20)
	maxLimit := getEnvAsInt64("MAX_SEARCH_LIMIT", 100)
	timeoutSeconds := getEnvAsInt64("GROUP_SERVICE_TIMEOUT_SECONDS", 5)

	return &Config{
		Env:               getEnv("APP_ENV", "development"),
		AppPort:           getEnv("SEARCH_SERVICE_PORT", getEnv("APP_PORT", "4000")),
		GroupServiceURL:   getEnv("GROUP_SERVICE_URL", "http://app-service:3001"),
		EnablePlayground:  enablePlayground,
		DefaultResultSize: defaultLimit,
		MaxResultSize:     maxLimit,
		HTTPClientTimeout: time.Duration(timeoutSeconds) * time.Second,
	}, nil
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		return value
	}
	return fallback
}

func getEnvAsBool(key string, fallback bool) bool {
	valStr, exists := os.LookupEnv(key)
	if !exists || valStr == "" {
		return fallback
	}
	val, err := strconv.ParseBool(valStr)
	if err != nil {
		return fallback
	}
	return val
}

func getEnvAsInt64(key string, fallback int64) int64 {
	valStr, exists := os.LookupEnv(key)
	if !exists || valStr == "" {
		return fallback
	}
	val, err := strconv.ParseInt(valStr, 10, 64)
	if err != nil {
		return fallback
	}
	return val
}

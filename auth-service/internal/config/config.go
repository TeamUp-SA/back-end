package config

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
    GoogleOAuthClientID     string
    GoogleOAuthClientSecret string
    OAuthRedirectURL        string
    SessionSecret           string
    DatabaseDSN             string
    JWTSecret               string
    UserServiceAddr         string
    AppServiceBaseURL       string
    AllowedOrigins          []string
}

func Load() *Config {
	// Load .env if present
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system env vars")
	}

    return &Config{
        GoogleOAuthClientID:     os.Getenv("GOOGLE_OAUTH_CLIENT_ID"),
        GoogleOAuthClientSecret: os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET"),
        OAuthRedirectURL:        os.Getenv("OAUTH_REDIRECT_URL"),
        SessionSecret:           os.Getenv("SESSION_SECRET"),
        DatabaseDSN:             os.Getenv("DATABASE_DSN"),
        JWTSecret:               os.Getenv("JWT_SECRET"),
        UserServiceAddr:         getenvDefault("USER_SERVICE_ADDR", "localhost:9091"),
        AppServiceBaseURL:       getenvDefault("APP_SERVICE_BASE_URL", "http://app-service:3001"),
        AllowedOrigins:          parseCSVEnv("AUTH_ALLOWED_ORIGINS", []string{"http://localhost:3000", "http://localhost:3001"}),
    }
}

func getenvDefault(k, def string) string {
    v := os.Getenv(k)
    if v == "" { return def }
    return v
}

func parseCSVEnv(key string, def []string) []string {
    raw := strings.TrimSpace(os.Getenv(key))
    if raw == "" {
        return def
    }
    parts := strings.Split(raw, ",")
    out := make([]string, 0, len(parts))
    for _, p := range parts {
        trimmed := strings.TrimSpace(p)
        if trimmed != "" {
            out = append(out, trimmed)
        }
    }
    if len(out) == 0 {
        return def
    }
    return out
}

package config

import (
	"log"
	"os"

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
        UserServiceAddr:         getenvDefault("USER_SERVICE_ADDR", "localhost:9091"),
    }
}

func getenvDefault(k, def string) string {
    v := os.Getenv(k)
    if v == "" { return def }
    return v
}

package config

import (
    "log"
    "os"

    "github.com/joho/godotenv"
)

type Config struct {
    DatabaseDSN string
    JWTSecret   string
    SessionSecret string
}

func Load() *Config {
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found, using system env vars")
    }
    return &Config{
        DatabaseDSN:  os.Getenv("DATABASE_DSN"),
        JWTSecret:    os.Getenv("JWT_SECRET"),
        SessionSecret: os.Getenv("SESSION_SECRET"),
    }
}


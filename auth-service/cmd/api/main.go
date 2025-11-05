package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	"auth-service/internal/auth"
	"auth-service/internal/config"
	h "auth-service/internal/http/handlers"
)

var ctx = context.Background()

func main() {
	// Load .env or environment variables
	cfg := config.Load()

	// Initialize Google OAuth credentials
	auth.InitGoogleOAuth()

	// Ensure JWT secret is available
	if cfg.JWTSecret == "" {
		log.Fatal("❌ JWT_SECRET is required")
	}

	// ✅ Initialize Redis client
	rdb := redis.NewClient(&redis.Options{
		Addr:     "redis:6379", 
		Password: "",          
		DB:       0,          
	})

	// Ping Redis
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf("❌ Failed to connect to Redis: %v", err)
	}
	log.Println("✅ Connected to Redis")

	// Initialize Gin router
	r := gin.Default()

	// Inject redis into handlers
	handler := h.NewHandler(rdb)

	// Public routes
	r.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "Auth service is running.")
	})

	// Auth routes
	r.GET("/login", handler.LoginHandler)
	r.GET("/auth/callback", handler.CallbackHandler)
	r.GET("/logout", handler.LogoutHandler)

	// Start the service
	if err := r.Run(":8082"); err != nil {
		log.Fatal(err)
	}
}

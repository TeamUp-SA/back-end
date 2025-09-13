package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	// internal packages (update the module path to your actual module name from go.mod)
	"auth-service/internal/auth"
	"auth-service/internal/config"
	h "auth-service/internal/http/handlers"
)

func main() {
	// Load configuration (.env is handled inside config.Load)
	cfg := config.Load()
	db.MustInitGorm()
	if err := db.Gorm().AutoMigrate(&auth.User{}); err != nil {
		log.Fatalf("auto-migrate: %v", err)
	}

	// Initialize Google OAuth (reads env vars internally)
	auth.InitGoogleOAuth()

	// Ensure JWT secret is set (for HS256). For RS256, switch to key files instead.
	if cfg.SessionSecret == "" && cfg.JWTSecret == "" {
		log.Fatal("JWT_SECRET is required when using JWT auth (or configure RS256 keys)")
	}

	// Gin router
	r := gin.Default()

	// Public routes
	r.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello! Go to /login to sign in with Google.")
	})
	r.GET("/login", h.LoginHandler)
	r.GET("/auth/callback", h.CallbackHandler)
	r.GET("/logout", h.LogoutHandler)

	// Start server
	if err := r.Run(":8082"); err != nil {
		log.Fatal(err)
	}
}

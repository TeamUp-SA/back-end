package main

import (
	"encoding/gob"
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"

	// internal packages (update the module path to your actual module name from go.mod)
	"user-service/internal/auth"
	"user-service/internal/config"
	h "user-service/internal/http/handlers"
	m "user-service/internal/http/middleware"
)

func main() {
	// Load configuration (.env is handled inside config.Load)
	cfg := config.Load()

	// Register types stored in sessions
	gob.Register(auth.GoogleUser{})

	// Initialize Google OAuth (reads env vars internally)
	auth.InitGoogleOAuth()

	// Gin router
	r := gin.Default()

	// Session store (cookie-based). In production, ensure HTTPS and strong secrets.
	if cfg.SessionSecret == "" {
		log.Fatal("SESSION_SECRET is required")
	}
	store := cookie.NewStore([]byte(cfg.SessionSecret))
	store.Options(sessions.Options{
		Path:     "/",
		MaxAge:   60 * 60 * 8, // 8 hours
		HttpOnly: true,
		Secure:   false, // set true when behind HTTPS/production
		SameSite: http.SameSiteLaxMode,
	})
	r.Use(sessions.Sessions(auth.SessionName, store))

	// Public routes
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello! Go to /login to sign in with Google.")
	})
	r.GET("/login", h.LoginHandler)
	r.GET("/auth/callback", h.CallbackHandler)
	r.GET("/logout", h.LogoutHandler)

	// Protected routes
	app := r.Group("/app")
	app.Use(m.AuthRequired())
	app.GET("/me", h.GetMe)

	// Start server
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}

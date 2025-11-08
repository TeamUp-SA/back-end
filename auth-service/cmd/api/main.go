package main

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	"auth-service/internal/auth"
	"auth-service/internal/config"
	"auth-service/internal/db"
	h "auth-service/internal/http/handlers"
)

var ctx = context.Background()

func main() {
	// Load .env or environment variables
	cfg := config.Load()

	// Initialize database connection
	db.MustInitGorm()

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

	origins := cfg.AllowedOrigins
	if len(origins) == 0 {
		origins = []string{"http://localhost:3000", "http://localhost:3001"}
	}
	allowAll := false
	originSet := make(map[string]struct{}, len(origins))
	for _, o := range origins {
		o = strings.TrimSpace(o)
		if o == "" {
			continue
		}
		if o == "*" {
			allowAll = true
			continue
		}
		originSet[strings.ToLower(o)] = struct{}{}
	}

	r.Use(func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		allowed := false
		if allowAll {
			allowed = true
		} else if origin != "" {
			_, allowed = originSet[strings.ToLower(origin)]
		}

		if allowed && origin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		}
		if allowAll && origin == "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		}
		if allowed || allowAll {
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Authorization, Content-Type, Accept")
			c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
			c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
			c.Writer.Header().Set("Vary", "Origin")
		}

		if c.Request.Method == http.MethodOptions {
			if allowed || allowAll {
				c.AbortWithStatus(http.StatusNoContent)
				return
			}
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		c.Next()
	})

	// Inject redis into handlers
	handler := h.NewHandler(rdb, cfg)

	// Public routes
	r.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "Auth service is running.")
	})

	// Auth routes
	r.POST("/register", handler.RegisterHandler)
	r.POST("/login", handler.LoginHandler)
	r.GET("/login/google", handler.LoginWithGoogleHandler)
	r.GET("/auth/callback", handler.CallbackHandler)
	r.GET("/logout", handler.LogoutHandler)

	if err := r.Run(":8082"); err != nil {
		log.Fatal(err)
	}
}

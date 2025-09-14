package main

import (
	"log"
	"net"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"

	"user-service/internal/config"
	"user-service/internal/db"
	grpcserver "user-service/internal/grpc"
	h "user-service/internal/http/handlers"
	m "user-service/internal/http/middleware"
	"user-service/internal/user"
	userv1 "user-service/pb/userv1"
)

func main() {
	// Load env and init DB
	_ = config.Load()
	db.MustInitGorm()
	if err := db.Gorm().AutoMigrate(&user.User{}); err != nil {
		log.Fatalf("auto-migrate: %v", err)
	}

	r := gin.Default()

	// Health/public
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "user-service up"})
	})

	// Authenticated routes
	auth := r.Group("/api").Use(m.JWTAuth())
	{
		auth.GET("/me", h.GetMe)
		auth.PUT("/profile", h.UpdateProfile)
	}

	// Start gRPC server in a goroutine
	go func() {
		addr := os.Getenv("USER_GRPC_ADDR")
		if addr == "" {
			addr = ":9091"
		}
		lis, err := net.Listen("tcp", addr)
		if err != nil {
			log.Fatalf("grpc listen: %v", err)
		}
		s := grpc.NewServer()
		userv1.RegisterUserServiceServer(s, grpcserver.New())
		log.Printf("user-service gRPC listening on %s", addr)
		if err := s.Serve(lis); err != nil {
			log.Fatalf("grpc serve: %v", err)
		}
	}()

	if err := r.Run(":8081"); err != nil { // user-service on 8081
		log.Fatal(err)
	}
}

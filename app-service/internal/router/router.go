package router

import (
	"fmt"

	"time"

	docs "app-service/docs"
	"app-service/internal/config"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
	"go.mongodb.org/mongo-driver/mongo"
)

type Router struct {
	g    *gin.Engine
	conf *config.Config
	deps *Dependencies
}

func NewRouter(g *gin.Engine, conf *config.Config) *Router {
	return &Router{g, conf, nil}
}

func (r *Router) Run(mongoDB *mongo.Database) {

	// CORS setting - allow browser apps on localhost with credentials
	corsConfig := cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:3001"},
		AllowMethods:     []string{"OPTIONS", "PATCH", "PUT", "GET", "POST", "DELETE"},
		AllowHeaders:     []string{"Content-Type", "Authorization", "X-Member-ID", "x-member-id", "X-Requested-With", "Accept"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type"},
		MaxAge:           12 * time.Hour,
		AllowCredentials: true,
	}

	r.g.Use(cors.New(corsConfig))

	r.g.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "OK",
		})
	})

	// Swagger setting
	docs.SwaggerInfo.BasePath = ""
	r.g.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// versioning
	v1 := r.g.Group("")

	// setup
	r.deps = NewDependencies(mongoDB, r.conf)

	// Add related path
	r.AddBulletinRouter(v1)
	r.AddGroupRouter(v1)
	r.AddMemberRouter(v1)

	err := r.g.Run(":" + r.conf.App.Port)
	if err != nil {
		panic(fmt.Sprintf("Failed to run the server : %v", err))
	}
}

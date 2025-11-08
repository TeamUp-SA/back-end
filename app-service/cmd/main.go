package main

import (
	"fmt"

	"app-service/internal/config"
	"app-service/internal/database"
	routes "app-service/internal/router"
	"github.com/gin-gonic/gin"
)

func main() {
	conf, err := config.LoadConfig()

	if err != nil {
		panic(fmt.Sprintf("Error loading config: %v", err))
	}

	mongoDB, err := database.InitMongoDatabase(&conf.Db)

	if err != nil {
		panic(fmt.Sprintf("Error connecting mongo: %v", err))
	}

	r := routes.NewRouter(gin.Default(), conf)

	r.Run(mongoDB)
}

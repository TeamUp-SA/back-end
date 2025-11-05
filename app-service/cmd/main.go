package main

import (
	"fmt"

	"github.com/Ntchah/TeamUp-application-service/internal/config"
	"github.com/Ntchah/TeamUp-application-service/internal/database"
	routes "github.com/Ntchah/TeamUp-application-service/internal/router"
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

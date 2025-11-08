package database

import (
	"context"
	"fmt"
	"log"

	"app-service/internal/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InitMongoDatabase(conf *config.DbConfig) (db *mongo.Database, err error) {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(conf.MongoURL))
	if err != nil {
		panic(fmt.Sprintf("Error connecting mongo: %v", err))
	}
	database := client.Database("TeamUp")
	log.Println("✅ successfully connect MongoDB")
	return database, nil
}

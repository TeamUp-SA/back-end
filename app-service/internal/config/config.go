package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	Port string
	Env  string
}

type DbConfig struct {
	MongoURL string
}

type KafkaConfig struct {
	Broker           string
	NotificationTopic string
}

type Config struct {
	App AppConfig
	Db  DbConfig
	Kafka KafkaConfig
}

func LoadConfig() (*Config, error) {
	dir, err := os.Getwd() // Capture both the directory and the error
	if err != nil {
		fmt.Println("Error getting working directory:", err)
	}
	fmt.Println("Current Working Directory:", dir)
	if os.Getenv("APP_ENV") == "" {
		err := godotenv.Load(".env")
		if err != nil {
			return nil, err
		}
	}

	appConfig := AppConfig{
		Env:  os.Getenv("APP_ENV"),
		Port: os.Getenv("APP_PORT"),
	}

	dbConfig := DbConfig{
		MongoURL: os.Getenv("MONGODB_URL"),
	}

	kafkaConfig := KafkaConfig{
		Broker:           getEnv("KAFKA_BROKER", "kafka:9092"),
		NotificationTopic: getEnv("KAFKA_TOPIC", "notifications"),
	}

	return &Config{
		App:   appConfig,
		Db:    dbConfig,
		Kafka: kafkaConfig,
	}, nil
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		return value
	}
	return fallback
}

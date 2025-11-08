package config

import (
	"os"
)

type Config struct {
	KafkaBroker string
	KafkaTopic  string
	GroupID     string
}

func LoadConfig() *Config {
	broker := getEnv("KAFKA_BROKER", "localhost:9092")
	topic := getEnv("KAFKA_TOPIC", "notifications")
	groupID := getEnv("KAFKA_GROUP_ID", "notification-group")

	return &Config{
		KafkaBroker: broker,
		KafkaTopic:  topic,
		GroupID:     groupID,
	}
}

func getEnv(key, defaultVal string) string {
	if val, exists := os.LookupEnv(key); exists {
		return val
	}
	return defaultVal
}

package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"notification-service/internal/notifier"

	"github.com/segmentio/kafka-go"
)

type Message struct {
	Type    string `json:"type"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

func main() {
	broker := os.Getenv("KAFKA_BROKER")
	if broker == "" {
		broker = "kafka:9092"
	}

	topic := os.Getenv("KAFKA_TOPIC")
	if topic == "" {
		topic = "notifications"
	}

	groupID := os.Getenv("KAFKA_GROUP_ID")
	if groupID == "" {
		groupID = "notification-group"
	}

	ctx := context.Background()

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{broker},
		Topic:    topic,
		GroupID:  groupID,
		MinBytes: 1,
		MaxBytes: 10e6,
	})
	defer reader.Close()

	log.Println("Notification service started. Listening for messages...")

	for {
		msg, err := reader.ReadMessage(ctx)
		if err != nil {
			log.Println("Error reading message (Kafka might not be ready yet):", err)
			time.Sleep(3 * time.Second)
			continue
		}

		var m Message
		if err := json.Unmarshal(msg.Value, &m); err != nil {
			log.Println("Invalid message format:", err)
			continue
		}

		switch m.Type {
		case "email":
			notifier.SendEmail(m.To, m.Subject, m.Message)
		default:
			log.Println("Unknown message type:", m.Type)
		}
	}
}

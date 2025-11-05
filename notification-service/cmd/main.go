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

		var reader *kafka.Reader
		for i := 0; i < 10; i++ {
			reader = kafka.NewReader(kafka.ReaderConfig{
				Brokers: []string{broker},
				Topic:   topic,
				GroupID: groupID,
			})

			_, err := reader.FetchMessage(ctx)
			if err == nil {
				log.Println("Connected to Kafka successfully")
				break
			}

			log.Printf("Kafka not ready (%v), retrying in 3s...\n", err)
			time.Sleep(3 * time.Second)
		}
		if reader == nil {
			log.Fatal("Failed to connect to Kafka after retries")
		}
		defer reader.Close()

		log.Println("Notification service started. Listening for messages...")

		for {
			msg, err := reader.ReadMessage(ctx)
			if err != nil {
				log.Println("Error reading message:", err)
				continue
			}

			var m Message
			if err := json.Unmarshal(msg.Value, &m); err != nil {
				log.Println("Invalid message format:", err)
				continue
			}

			switch m.Type {
			case "email":
				notifier.SendEmail(m.To, m.Message)
			case "push":
				notifier.SendPush("This is Kafka message!", m.Message)
			default:
				log.Println("Unknown message type:", m.Type)
			}
		}
	}

// {"type": "email", "to": "dalai2547@gmail.com", "message": "Hello from Kafka email test!"}
// {"type": "push", "to": "test_user", "message": "You got a new Kafka push message!"}

package main

import (
	"log"
	"notification-service/internal/kafka"
)

func main() {
    consumer, err := kafka.NewConsumer("localhost:9092", "notifications", "notification-group")
    if err != nil {
        log.Fatalf("failed to create consumer: %v", err)
    }

    log.Println("Notification service started. Listening for messages...")
    consumer.ConsumeMessages()
}

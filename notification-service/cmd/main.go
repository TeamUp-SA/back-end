package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http" // Required for the server
	"os"
	"time"

	"notification-service/internal/notifier" // Assuming this package exists

	"github.com/segmentio/kafka-go"
)

type Message struct {
	Type    string `json:"type"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

func main() {
	// --- Configuration Setup ---
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

	// --- 1. Start Kafka Consumer in a Goroutine ---
	go startKafkaConsumer(ctx, broker, topic, groupID)

	// --- 2. Start Network Listeners (for readiness checks and API) ---
    // Start listener on 8083 (HTTP API port)
	go startListener(":8083", "HTTP API") 
    
    // Start listener on 9999 (Internal RPC/Readiness Check port)
	go startListener(":9999", "Internal RPC/Readiness")

	// --- 3. Block Main Thread ---
	log.Println("Notification service main thread running.")
	// select{} blocks indefinitely, keeping all goroutines alive
	select {} 
}

// startKafkaConsumer contains your original Kafka logic
func startKafkaConsumer(ctx context.Context, broker, topic, groupID string) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{broker},
		Topic:    topic,
		GroupID:  groupID,
		MinBytes: 1,
		MaxBytes: 10e6,
	})
	defer reader.Close()

	log.Println("✅ Kafka consumer started. Listening for messages on topic:", topic)

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
			// Assuming notifier.SendEmail exists and works
			notifier.SendEmail(m.To, m.Subject, m.Message)
		default:
			log.Println("Unknown message type:", m.Type)
		}
	}
}

func startListener(port, name string) {
	// A basic handler for health checks
    handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("OK"))
    })
    
    log.Printf("✅ %s listener starting on port %s...", name, port)

	if err := http.ListenAndServe(port, handler); err != nil {
		// Use Fatalf to crash the application if the port is unavailable
		log.Fatalf("Fatal error starting %s server on port %s: %v", name, port, err)
	}
}
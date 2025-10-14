package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"notification-service/internal/notifier"

	"github.com/segmentio/kafka-go"
)

type NotificationMessage struct {
    Type    string `json:"type"`
    To      string `json:"to"`
    Message string `json:"message"`
}

type Consumer struct {
    reader *kafka.Reader
}

func NewConsumer(broker, topic, groupID string) (*Consumer, error) {
    r := kafka.NewReader(kafka.ReaderConfig{
        Brokers: []string{broker},
        Topic:   topic,
        GroupID: groupID,
    })
    return &Consumer{reader: r}, nil
}

func (c *Consumer) ConsumeMessages() {
    for {
        m, err := c.reader.ReadMessage(context.Background())
        if err != nil {
            log.Printf("error reading message: %v", err)
            continue
        }

        var msg NotificationMessage
        if err := json.Unmarshal(m.Value, &msg); err != nil {
            log.Printf("invalid message format: %v", err)
            continue
        }

        log.Printf("Received notification: %+v", msg)

        switch msg.Type {
        case "email":
            notifier.SendEmail(msg.To, msg.Message)
        case "sms":
            notifier.SendSMS(msg.To, msg.Message)
        case "push":
            notifier.SendPush(msg.To, msg.Message)
        default:
            fmt.Printf("unknown type: %s\n", msg.Type)
        }
    }
}

package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/segmentio/kafka-go"
)

const (
	// NotificationTypeEmail represents an email notification type.
	NotificationTypeEmail = "email"
)

// NotificationMessage represents the payload pushed to Kafka.
type NotificationMessage struct {
	Type    string `json:"type"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

// Producer abstracts the notification producer behavior.
type Producer interface {
	Publish(ctx context.Context, messages []NotificationMessage) error
}

// NotificationProducer wraps a kafka writer to publish notification payloads.
type NotificationProducer struct {
	writer *kafka.Writer
}

func NewNotificationProducer(broker, topic string) (*NotificationProducer, error) {
	if broker == "" {
		return nil, errors.New("kafka producer: broker is required")
	}
	if topic == "" {
		return nil, errors.New("kafka producer: topic is required")
	}

	writer := &kafka.Writer{
		Addr:     kafka.TCP(broker),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}

	return &NotificationProducer{
		writer: writer,
	}, nil
}

func (p *NotificationProducer) Publish(ctx context.Context, messages []NotificationMessage) error {
	if p == nil || p.writer == nil {
		return errors.New("kafka producer: writer is not initialized")
	}
	if len(messages) == 0 {
		return nil
	}

	kafkaMessages := make([]kafka.Message, 0, len(messages))
	for _, message := range messages {
		payload, err := json.Marshal(message)
		if err != nil {
			return fmt.Errorf("kafka producer: marshal message: %w", err)
		}

		kafkaMessages = append(kafkaMessages, kafka.Message{
			Key:   []byte(message.To),
			Value: payload,
		})
	}

	if err := p.writer.WriteMessages(ctx, kafkaMessages...); err != nil {
		return fmt.Errorf("kafka producer: write messages: %w", err)
	}

	return nil
}



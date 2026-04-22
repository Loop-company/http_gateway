package kafka

import (
	"context"
	"encoding/json"

	"github.com/segmentio/kafka-go"
)

const loginCommandTopic = "auth.command.login"

type AuthProducer struct {
	writer *kafka.Writer
}

func NewAuthProducer(brokers []string) *AuthProducer {
	return &AuthProducer{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(brokers...),
			Topic:    loginCommandTopic,
			Balancer: &kafka.LeastBytes{},
			// Async:    false,
		},
	}
}

func (p *AuthProducer) Close() error {
	return p.writer.Close()
}

func (p *AuthProducer) SendLogin(ctx context.Context, cmd LoginCommand) error {
	data, err := json.Marshal(cmd)
	if err != nil {
		return err
	}

	return p.writer.WriteMessages(ctx, kafka.Message{
		Value: data,
	})
}

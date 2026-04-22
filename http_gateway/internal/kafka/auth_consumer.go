package kafka

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/segmentio/kafka-go"
)

const (
	loginResponseTopic = "auth.event.token_generated"
	groupID            = "gateway-login-responses"
)

type AuthConsumer struct {
	reader    *kafka.Reader
	responses map[string]chan LoginResponse
	mu        sync.RWMutex
}

func NewAuthConsumer(brokers []string) *AuthConsumer {
	return &AuthConsumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:  brokers,
			GroupID:  groupID,
			Topic:    loginResponseTopic,
			MinBytes: 10e3, // 10 KiB
			MaxBytes: 10e6, // 10 MiB
		}),
		responses: make(map[string]chan LoginResponse),
	}
}

func (c *AuthConsumer) Close() error {
	return c.reader.Close()
}

func (c *AuthConsumer) RegisterRequest(requestID string) chan LoginResponse {
	ch := make(chan LoginResponse, 1)

	c.mu.Lock()
	c.responses[requestID] = ch
	c.mu.Unlock()

	return ch
}

func (c *AuthConsumer) UnregisterRequest(requestID string) {
	c.mu.Lock()
	delete(c.responses, requestID)
	c.mu.Unlock()
}

func (c *AuthConsumer) Run(ctx context.Context) {
	for {
		msg, err := c.reader.ReadMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return // graceful shutdown
			}
			time.Sleep(2 * time.Second)
			continue
		}

		var resp LoginResponse
		if err := json.Unmarshal(msg.Value, &resp); err != nil {
			continue
		}

		c.mu.Lock()
		if ch, ok := c.responses[resp.RequestID]; ok {
			ch <- resp
			close(ch)
			delete(c.responses, resp.RequestID)
		}
		c.mu.Unlock()
	}
}

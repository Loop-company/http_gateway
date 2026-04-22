package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/Loop-company/http_gateway/internal/kafka"
)

type AuthHandler struct {
	producer *kafka.AuthProducer
	consumer *kafka.AuthConsumer
}

func NewAuthHandler(producer *kafka.AuthProducer, consumer *kafka.AuthConsumer) *AuthHandler {
	return &AuthHandler{
		producer: producer,
		consumer: consumer,
	}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email"    binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	requestID := uuid.New().String()

	cmd := kafka.LoginCommand{
		RequestID: requestID,
		Email:     req.Email,
		Password:  req.Password,
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 12*time.Second)
	defer cancel()

	// Регистрируем ожидание ответа
	responseCh := h.consumer.RegisterRequest(requestID)
	defer h.consumer.UnregisterRequest(requestID)

	// Отправляем команду в Kafka
	if err := h.producer.SendLogin(ctx, cmd); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to send login command"})
		return
	}

	select {
	case resp := <-responseCh:
		if resp.Error != "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": resp.Error})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"access_token":  resp.AccessToken,
			"refresh_token": resp.RefreshToken,
		})

	case <-ctx.Done():
		c.JSON(http.StatusGatewayTimeout, gin.H{"error": "timeout waiting for auth service"})
	}
}

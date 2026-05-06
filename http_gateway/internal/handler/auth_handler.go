package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/Loop-company/http_gateway/internal/authclient"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	client *authclient.Client
}

func NewAuthHandler(client *authclient.Client) *AuthHandler {
	return &AuthHandler{client: client}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 12*time.Second)
	defer cancel()

	userAgent := c.GetHeader("User-Agent")
	ip := c.ClientIP()

	resp, statusCode, err := h.client.Login(ctx, req.Email, req.Password, userAgent, ip)
	if err != nil {
		if statusCode == http.StatusUnauthorized {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}

		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"guid":               resp.GUID,
		"access_token":       resp.AccessToken,
		"refresh_token":      resp.RefreshToken,
		"access_expires_at":  resp.AccessExpiresAt,
		"refresh_expires_at": resp.RefreshExpiresAt,
	})
}

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

func (h *AuthHandler) Register(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	msg, err := h.client.Register(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"message": msg})
}

func (h *AuthHandler) Verify(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
		Code  string `json:"code" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	guid, err := h.client.Verify(c.Request.Context(), req.Email, req.Code)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid code"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"guid": guid})
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	userGUID, _ := c.Get("user_guid")
	sessionID, _ := c.Get("session_id")

	userAgent := c.GetHeader("User-Agent")
	ip := c.ClientIP()

	resp, err := h.client.Refresh(c.Request.Context(), req.RefreshToken, userAgent, ip, userGUID.(string), sessionID.(string))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	userGUID, _ := c.Get("user_guid")

	err := h.client.Logout(c.Request.Context(), userGUID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
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

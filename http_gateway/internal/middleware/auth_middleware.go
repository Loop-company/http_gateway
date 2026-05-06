package middleware

import (
	"net/http"
	"strings"

	"github.com/Loop-company/http_gateway/internal/authclient"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(client *authclient.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			// Check cookie
			cookie, err := c.Cookie("access_token")
			if err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
				return
			}
			authHeader = "Bearer " + cookie
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid auth header"})
			return
		}

		accessToken := parts[1]

		guid, sessionID, err := client.ValidateToken(c.Request.Context(), accessToken)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		c.Set("user_guid", guid)
		c.Set("session_id", sessionID)
		c.Next()
	}
}

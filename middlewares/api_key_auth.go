package middlewares

import (
	"net/http"
	"strings"

	"github.com/BlenDMinh/dutgrad-server/controllers"
	"github.com/gin-gonic/gin"
)

func RequireApiKey() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing or invalid Authorization header"})
			c.Abort()
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")

		apiKey, err := controllers.VerifySpaceAPIKey(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired API Key token"})
			c.Abort()
			return
		}

		c.Set("apiKey", apiKey)
		c.Next()
	}
}

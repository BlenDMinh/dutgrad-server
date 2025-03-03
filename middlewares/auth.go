package middlewares

import (
	"net/http"

	"github.com/BlenDMinh/dutgrad-server/helpers"
	"github.com/BlenDMinh/dutgrad-server/models"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get token from Authorization header
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" || len(authHeader) < 8 || authHeader[:7] != "Bearer " {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, models.NewErrorResponse(http.StatusUnauthorized, "Unauthorized", nil))
			return
		}

		tokenString := authHeader[7:]

		// Verify token
		userID, err := helpers.VerifyJWTToken(tokenString)
		if err != nil {
			errMsg := err.Error()
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, models.NewErrorResponse(http.StatusUnauthorized, "Invalid or expired token", &errMsg))
			return
		}

		// Set userID in context
		ctx.Set("userID", userID)
		ctx.Next()
	}
}

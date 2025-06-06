package middlewares

import (
	"github.com/BlenDMinh/dutgrad-server/services"
	"github.com/gin-gonic/gin"
)

func ChatRateLimiter(userService services.UserService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userIDVal, exists := ctx.Get("user_id")
		if !exists || userIDVal == nil {
			ctx.Next()
			return
		}

		userID := userIDVal.(uint)

		if userService.IsRateLimited(userID) {
			ctx.JSON(429, gin.H{
				"error": "Rate limit exceeded. Please try again later.",
			})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

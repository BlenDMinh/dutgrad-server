package controllers

import (
	"net/http"
	"strconv"

	"github.com/BlenDMinh/dutgrad-server/models"
	"github.com/gin-gonic/gin"
)

// ExtractID extracts an ID from either URL parameters or context
// If paramName is "user_id", extracts user ID from context
// Otherwise extracts from URL parameters
// Returns the ID and true if successful, 0 and false otherwise (also handles error response)
func ExtractID(ctx *gin.Context, paramName string) (uint, bool) {
	// For user ID from context
	if paramName == "user_id" {
		userID, exists := ctx.Get("user_id")
		if !exists {
			ctx.JSON(http.StatusUnauthorized, models.NewErrorResponse(
				http.StatusUnauthorized,
				"User ID not found in context",
				nil,
			))
			return 0, false
		}
		return userID.(uint), true
	}

	// For ID from URL parameter
	idStr := ctx.Param(paramName)
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		errMsg := err.Error()
		ctx.JSON(http.StatusBadRequest, models.NewErrorResponse(
			http.StatusBadRequest,
			"Invalid ID parameter",
			&errMsg,
		))
		return 0, false
	}
	return uint(id), true
}

// HandleBindJSON binds and validates JSON request
func HandleBindJSON[T any](ctx *gin.Context, req *T) bool {
	if err := ctx.ShouldBindJSON(req); err != nil {
		errMsg := err.Error()
		ctx.JSON(http.StatusBadRequest, models.NewErrorResponse(
			http.StatusBadRequest,
			"Invalid request",
			&errMsg,
		))
		return false
	}
	return true
}

// HandleSuccess handles successful responses
func HandleSuccess(ctx *gin.Context, message string, data interface{}) {
	ctx.JSON(http.StatusOK, models.NewSuccessResponse(
		http.StatusOK,
		message,
		data,
	))
}

// HandleCreated handles creation success responses
func HandleCreated(ctx *gin.Context, message string, data interface{}) {
	ctx.JSON(http.StatusCreated, models.NewSuccessResponse(
		http.StatusCreated,
		message,
		data,
	))
}

// HandleError handles error responses
func HandleError(ctx *gin.Context, statusCode int, message string, err error) {
	var errMsg *string
	if err != nil {
		errString := err.Error()
		errMsg = &errString
	}

	ctx.JSON(statusCode, models.NewErrorResponse(
		statusCode,
		message,
		errMsg,
	))
}

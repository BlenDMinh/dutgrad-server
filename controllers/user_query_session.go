package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/BlenDMinh/dutgrad-server/databases/entities"
	"github.com/BlenDMinh/dutgrad-server/models"
	"github.com/BlenDMinh/dutgrad-server/models/dtos"
	"github.com/BlenDMinh/dutgrad-server/services"
	"github.com/gin-gonic/gin"
)

type UserQuerySessionController struct {
	CrudController[entities.UserQuerySession, uint]
}

func NewUserQuerySessionController() *UserQuerySessionController {
	return &UserQuerySessionController{
		CrudController: *NewCrudController(services.NewUserQuerySessionService()),
	}
}

func (c *UserQuerySessionController) BeginChatSession(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusInternalServerError, models.NewErrorResponse(
			http.StatusInternalServerError,
			"User ID not found in context",
			nil,
		))
		return
	}

	var req dtos.BeginChatSessionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		errMsg := err.Error()
		ctx.JSON(http.StatusBadRequest, models.NewErrorResponse(
			http.StatusBadRequest,
			"Invalid request",
			&errMsg,
		))
		return
	}

	var userIDNum uint = userID.(uint)

	session := &entities.UserQuerySession{
		UserID:  &userIDNum,
		SpaceID: req.SpaceID,
	}

	session, err := c.service.Create(session)
	if err != nil {
		errMsg := err.Error()
		ctx.JSON(500, models.NewErrorResponse(
			http.StatusInternalServerError,
			"Failed to create session",
			&errMsg,
		))
		return
	}

	ctx.JSON(http.StatusOK, models.NewSuccessResponse(
		http.StatusOK,
		"Session created successfully",
		session,
	))
}

func (c *UserQuerySessionController) GetMyChatSessions(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, models.NewErrorResponse(
			http.StatusUnauthorized,
			"User ID not found in context",
			nil,
		))
		return
	}

	sessions, err := c.service.(*services.UserQuerySessionService).GetChatSessionsByUserID(userID.(uint))
	if err != nil {
		errMsg := err.Error()
		ctx.JSON(http.StatusInternalServerError, models.NewErrorResponse(
			http.StatusInternalServerError,
			"Failed to fetch sessions",
			&errMsg,
		))
		return
	}

	ctx.JSON(http.StatusOK, models.NewSuccessResponse(
		http.StatusOK,
		"Fetched sessions successfully",
		sessions,
	))
}

func (c *UserQuerySessionController) CountMyChatSessions(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.Status(http.StatusUnauthorized)
		return
	}

	count, err := c.service.(*services.UserQuerySessionService).CountChatSessionsByUserID(userID.(uint))
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.Header("X-Total-Count", fmt.Sprintf("%d", count))
	ctx.Status(http.StatusOK)
}

func (c *UserQuerySessionController) GetTempMessageByID(ctx *gin.Context) {
	sessionID, exists := ctx.Params.Get("id")
	if !exists {
		ctx.JSON(http.StatusBadRequest, models.NewErrorResponse(
			http.StatusBadRequest,
			"Invalid session ID",
			nil,
		))
		return
	}

	sessionIDNum, err := strconv.ParseUint(sessionID, 10, 32)

	if err != nil {
		errMsg := err.Error()
		ctx.JSON(http.StatusBadRequest, models.NewErrorResponse(
			http.StatusBadRequest,
			"Invalid session ID",
			&errMsg,
		))
		return
	}

	tempMessage, err := c.service.(*services.UserQuerySessionService).GetTempMessageByID(uint(sessionIDNum))
	if err != nil {
		errMsg := err.Error()
		ctx.JSON(http.StatusInternalServerError, models.NewErrorResponse(
			http.StatusInternalServerError,
			"Failed to fetch temp message",
			&errMsg,
		))
		return
	}

	ctx.JSON(http.StatusOK, models.NewSuccessResponse(
		http.StatusOK,
		"Fetched temp message successfully",
		tempMessage,
	))
}

func (c *UserQuerySessionController) GetChatHistory(ctx *gin.Context) {
	sessionID, exists := ctx.Params.Get("id")
	if !exists {
		ctx.JSON(http.StatusBadRequest, models.NewErrorResponse(
			http.StatusBadRequest,
			"Invalid session ID",
			nil,
		))
		return
	}

	sessionIDNum, err := strconv.ParseUint(sessionID, 10, 32)
	if err != nil {
		errMsg := err.Error()
		ctx.JSON(http.StatusBadRequest, models.NewErrorResponse(
			http.StatusBadRequest,
			"Invalid session ID",
			&errMsg,
		))
		return
	}

	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, models.NewErrorResponse(
			http.StatusUnauthorized,
			"User ID not found in context",
			nil,
		))
		return
	}

	history, err := c.service.(*services.UserQuerySessionService).GetChatHistoryBySessionID(uint(sessionIDNum), userID.(uint))
	if err != nil {
		errMsg := err.Error()
		ctx.JSON(http.StatusInternalServerError, models.NewErrorResponse(
			http.StatusInternalServerError,
			"Failed to fetch chat history",
			&errMsg,
		))
		return
	}

	ctx.JSON(http.StatusOK, models.NewSuccessResponse(
		http.StatusOK,
		"Fetched chat history successfully",
		history,
	))
}

func (c *UserQuerySessionController) ClearChatHistory(ctx *gin.Context) {
	sessionID, exists := ctx.Params.Get("id")
	if !exists {
		ctx.JSON(http.StatusBadRequest, models.NewErrorResponse(
			http.StatusBadRequest,
			"Invalid session ID",
			nil,
		))
		return
	}

	sessionIDNum, err := strconv.ParseUint(sessionID, 10, 32)
	if err != nil {
		errMsg := err.Error()
		ctx.JSON(http.StatusBadRequest, models.NewErrorResponse(
			http.StatusBadRequest,
			"Invalid session ID",
			&errMsg,
		))
		return
	}

	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, models.NewErrorResponse(
			http.StatusUnauthorized,
			"User ID not found in context",
			nil,
		))
		return
	}

	err = c.service.(*services.UserQuerySessionService).ClearChatHistoryBySessionID(uint(sessionIDNum), userID.(uint))
	if err != nil {
		errMsg := err.Error()
		ctx.JSON(http.StatusInternalServerError, models.NewErrorResponse(
			http.StatusInternalServerError,
			"Failed to clear chat history",
			&errMsg,
		))
		return
	}

	ctx.JSON(http.StatusOK, models.NewSuccessResponse(
		http.StatusOK,
		"Chat history cleared successfully",
		nil,
	))
}

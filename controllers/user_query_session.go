package controllers

import (
	"fmt"
	"net/http"

	"github.com/BlenDMinh/dutgrad-server/databases/entities"
	"github.com/BlenDMinh/dutgrad-server/models/dtos"
	"github.com/BlenDMinh/dutgrad-server/services"
	"github.com/gin-gonic/gin"
)

type UserQuerySessionController struct {
	CrudController[entities.UserQuerySession, uint]
	service services.UserQuerySessionService
}

func NewUserQuerySessionController(
	service services.UserQuerySessionService,
) *UserQuerySessionController {
	crudController := NewCrudController(service)
	return &UserQuerySessionController{
		CrudController: *crudController,
		service:        service,
	}
}

func (c *UserQuerySessionController) BeginChatSession(ctx *gin.Context) {
	userID, ok := ExtractID(ctx, "user_id")
	if !ok {
		return
	}

	var req dtos.BeginChatSessionRequest
	if !HandleBindJSON(ctx, &req) {
		return
	}

	session := &entities.UserQuerySession{
		UserID:  &userID,
		SpaceID: req.SpaceID,
	}

	session, err := c.service.Create(session)
	if err != nil {
		HandleError(ctx, http.StatusInternalServerError, "Failed to create session", err)
		return
	}

	HandleSuccess(ctx, "Session created successfully", session)
}

func (c *UserQuerySessionController) GetMyChatSessions(ctx *gin.Context) {
	userID, ok := ExtractID(ctx, "user_id")
	if !ok {
		return
	}

	sessions, err := c.service.GetChatSessionsByUserID(userID)
	if err != nil {
		HandleError(ctx, http.StatusInternalServerError, "Failed to fetch sessions", err)
		return
	}

	HandleSuccess(ctx, "Fetched sessions successfully", sessions)
}

func (c *UserQuerySessionController) CountMyChatSessions(ctx *gin.Context) {
	userID, ok := ExtractID(ctx, "user_id")
	if !ok {
		return
	}

	count, err := c.service.CountChatSessionsByUserID(userID)
	if err != nil {
		HandleError(ctx, http.StatusInternalServerError, "Failed to count chat sessions", err)
		return
	}

	ctx.Header("X-Total-Count", fmt.Sprintf("%d", count))
	ctx.Status(http.StatusOK)
}

func (c *UserQuerySessionController) GetTempMessageByID(ctx *gin.Context) {
	sessionID, ok := ExtractID(ctx, "id")
	if !ok {
		return
	}

	tempMessage, err := c.service.GetTempMessageByID(sessionID)
	if err != nil {
		HandleError(ctx, http.StatusInternalServerError, "Failed to fetch temp message", err)
		return
	}

	HandleSuccess(ctx, "Fetched temp message successfully", tempMessage)
}

func (c *UserQuerySessionController) GetChatHistory(ctx *gin.Context) {
	sessionID, ok := ExtractID(ctx, "id")
	if !ok {
		return
	}

	userID, ok := ExtractID(ctx, "user_id")
	if !ok {
		return
	}

	history, err := c.service.GetChatHistoryBySessionID(sessionID, userID)
	if err != nil {
		HandleError(ctx, http.StatusInternalServerError, "Failed to fetch chat history", err)
		return
	}

	HandleSuccess(ctx, "Fetched chat history successfully", history)
}

func (c *UserQuerySessionController) ClearChatHistory(ctx *gin.Context) {
	sessionID, ok := ExtractID(ctx, "id")
	if !ok {
		return
	}

	userID, ok := ExtractID(ctx, "user_id")
	if !ok {
		return
	}

	err := c.service.ClearChatHistoryBySessionID(sessionID, userID)
	if err != nil {
		HandleError(ctx, http.StatusInternalServerError, "Failed to clear chat history", err)
		return
	}

	HandleSuccess(ctx, "Chat history and session cleared successfully", nil)
}

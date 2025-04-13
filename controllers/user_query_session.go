package controllers

import (
	"net/http"

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

	session := &entities.UserQuerySession{
		UserID:  userID.(uint),
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

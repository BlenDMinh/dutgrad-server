package controllers

import (
	"net/http"

	"github.com/BlenDMinh/dutgrad-server/databases/entities"
	"github.com/BlenDMinh/dutgrad-server/models"
	"github.com/BlenDMinh/dutgrad-server/models/dtos"
	"github.com/BlenDMinh/dutgrad-server/services"
	"github.com/gin-gonic/gin"
)

type UserQueryController struct {
	CrudController[entities.UserQuery, uint]
}

func NewUserQueryController() *UserQueryController {
	return &UserQueryController{
		CrudController: *NewCrudController(services.NewUserQueryService()),
	}
}

func (c *UserQueryController) Ask(ctx *gin.Context) {
	_, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusInternalServerError, models.NewErrorResponse(
			http.StatusInternalServerError,
			"User ID not found in context",
			nil,
		))
		return
	}

	var req dtos.AskRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		errMsg := err.Error()
		ctx.JSON(http.StatusBadRequest, models.NewErrorResponse(
			http.StatusBadRequest,
			"Invalid request",
			&errMsg,
		))
		return
	}

	sessionService := services.NewUserQuerySessionService()

	// Check if session exists
	session, err := sessionService.GetById(req.QuerySessionID)
	if err != nil {
		errMsg := err.Error()
		ctx.JSON(http.StatusNotFound, models.NewErrorResponse(
			http.StatusNotFound,
			"Session not found",
			&errMsg,
		))
		return
	}

	ragService := services.NewRAGServerService()
	answer, err := ragService.Chat(req.QuerySessionID, session.ID, req.Query)
	if err != nil {
		errMsg := err.Error()
		ctx.JSON(http.StatusInternalServerError, models.NewErrorResponse(
			http.StatusInternalServerError,
			"Failed to get answer",
			&errMsg,
		))
		return
	}

	query := &entities.UserQuery{
		QuerySessionID: session.ID,
		Query:          req.Query,
	}

	query, err = c.service.Create(query)
	if err != nil {
		errMsg := err.Error()
		ctx.JSON(http.StatusInternalServerError, models.NewErrorResponse(
			http.StatusInternalServerError,
			"Failed to save query",
			&errMsg,
		))
		return
	}

	ctx.JSON(http.StatusOK, models.NewSuccessResponse(
		http.StatusOK,
		"Answer retrieved successfully",
		&gin.H{
			"answer": answer,
			"query":  query,
		},
	))
}

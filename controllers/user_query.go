package controllers

import (
	"net/http"

	"github.com/BlenDMinh/dutgrad-server/databases/entities"
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
	_, ok := ExtractID(ctx, "user_id")
	if !ok {
		return
	}

	var req dtos.AskRequest
	if !HandleBindJSON(ctx, &req) {
		return
	}

	sessionService := services.NewUserQuerySessionService()

	session, err := sessionService.GetById(req.QuerySessionID)
	if err != nil {
		HandleError(ctx, http.StatusNotFound, "Session not found", err)
		return
	}
	ragService := services.NewRAGServerService()
	answer, err := ragService.Chat(req.QuerySessionID, session.SpaceID, req.Query)
	if err != nil {
		HandleError(ctx, http.StatusInternalServerError, "Failed to get answer", err)
		return
	}

	query := &entities.UserQuery{
		QuerySessionID: session.ID,
		Query:          req.Query,
	}

	query, err = c.service.Create(query)
	query.UserQuerySession = *session

	if err != nil {
		HandleError(ctx, http.StatusInternalServerError, "Failed to save query", err)
		return
	}

	HandleSuccess(ctx, "Answer retrieved successfully", gin.H{
		"answer": answer,
		"query":  query,
	})
}

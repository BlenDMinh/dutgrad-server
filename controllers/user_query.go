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
	service services.UserQueryService
}

func NewUserQueryController(
	service services.UserQueryService,
) *UserQueryController {
	crudController := NewCrudController(service)
	return &UserQueryController{
		CrudController: *crudController,
		service:        service,
	}
}

func (c *UserQueryController) Ask(ctx *gin.Context) {
	userID, ok := ExtractID(ctx, "user_id")
	if !ok {
		return
	}
	var req dtos.AskRequest
	if !HandleBindJSON(ctx, &req) {
		return
	}

	if len(req.Query) > 1024 {
		req.Query = req.Query[:1024]
	}

	userService := services.NewUserService()
	tierUsage, err := userService.GetUserTierUsage(userID)
	if err != nil {
		HandleError(ctx, http.StatusInternalServerError, "Failed to check user tier usage", err)
		return
	}

	if tierUsage.Usage.ChatUsageDaily >= int64(tierUsage.Tier.QueryLimit) {
		HandleError(ctx, http.StatusTooManyRequests, "You have reached your daily chat limit. Please try again tomorrow or upgrade your plan.", nil)
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

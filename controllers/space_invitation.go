package controllers

import (
	"net/http"

	"github.com/BlenDMinh/dutgrad-server/databases/entities"
	"github.com/BlenDMinh/dutgrad-server/services"
	"github.com/gin-gonic/gin"
)

type SpaceInvitationController struct {
	CrudController[entities.SpaceInvitation, uint]
	service services.SpaceInvitationService
}

func NewSpaceInvitationController(
	service services.SpaceInvitationService,
) *SpaceInvitationController {
	crudController := NewCrudController(service)
	return &SpaceInvitationController{
		CrudController: *crudController,
		service:        service,
	}
}

func (c *SpaceInvitationController) AcceptInvitation(ctx *gin.Context) {
	userId, ok := ExtractID(ctx, "user_id")
	if !ok {
		return
	}

	invitationId, ok := ExtractID(ctx, "id")
	if !ok {
		return
	}

	err := c.service.AcceptInvitation(invitationId, userId)
	if err != nil {
		HandleError(ctx, http.StatusInternalServerError, "Failed to accept invitation", err)
		return
	}

	HandleSuccess(ctx, "Invitation accepted successfully", gin.H{"ok": "yes"})
}

func (c *SpaceInvitationController) RejectInvitation(ctx *gin.Context) {
	userId, ok := ExtractID(ctx, "user_id")
	if !ok {
		return
	}

	invitationId, ok := ExtractID(ctx, "id")
	if !ok {
		return
	}

	err := c.service.RejectInvitation(invitationId, userId)
	if err != nil {
		HandleError(ctx, http.StatusInternalServerError, "Failed to reject invitation", err)
		return
	}

	HandleSuccess(ctx, "Invitation rejected successfully", gin.H{"ok": "yes"})
}

func (c *SpaceInvitationController) GetInvitationCount(ctx *gin.Context) {
	userID, ok := ExtractID(ctx, "user_id")
	if !ok {
		return
	}

	count, err := c.service.CountInvitationByUserID(userID)
	if err != nil {
		HandleError(ctx, http.StatusInternalServerError, "Failed to get invitation count", err)
		return
	}

	HandleSuccess(ctx, "Invitation count retrieved successfully", gin.H{"count": count})
}

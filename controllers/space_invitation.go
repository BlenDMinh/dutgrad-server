package controllers

import (
	"net/http"
	"strconv"

	"github.com/BlenDMinh/dutgrad-server/databases/entities"
	"github.com/BlenDMinh/dutgrad-server/models"
	"github.com/BlenDMinh/dutgrad-server/services"
	"github.com/gin-gonic/gin"
)

type SpaceInvitationController struct {
	CrudController[entities.SpaceInvitation, uint]
}

func NewSpaceInvitationController() *SpaceInvitationController {
	return &SpaceInvitationController{
		CrudController: *NewCrudController(services.NewSpaceInvitationService()),
	}
}

func (c *SpaceInvitationController) AcceptInvitation(ctx *gin.Context) {
	userIdValue, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusInternalServerError, models.NewErrorResponse(http.StatusInternalServerError, "User ID not found in context", nil))
		return
	}
	userId := userIdValue.(uint)

	invitationIdParam := ctx.Param("id")
	invitationId, err := strconv.ParseUint(invitationIdParam, 10, 64)
	if err != nil {
		errMsg := err.Error()
		ctx.JSON(http.StatusBadRequest, models.NewErrorResponse(http.StatusBadRequest, "Invalid invitation ID", &errMsg))
		return
	}

	err = c.service.(*services.SpaceInvitationService).AcceptInvitation(uint(invitationId), userId)
	if err != nil {
		errMsg := err.Error()
		ctx.JSON(http.StatusInternalServerError, models.NewErrorResponse(http.StatusInternalServerError, "Failed to accept invitation", &errMsg))
		return
	}

	ctx.JSON(http.StatusOK, models.NewSuccessResponse(
		http.StatusOK,
		"Invitation accepted successfully",
		gin.H{"ok": "yes"},
	))
}

func (c *SpaceInvitationController) RejectInvitation(ctx *gin.Context) {
	userIdValue, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusInternalServerError, models.NewErrorResponse(http.StatusInternalServerError, "User ID not found in context", nil))
		return
	}
	userId := userIdValue.(uint)

	invitationIdParam := ctx.Param("id")
	invitationId, err := strconv.ParseUint(invitationIdParam, 10, 64)
	if err != nil {
		errMsg := err.Error()
		ctx.JSON(http.StatusBadRequest, models.NewErrorResponse(http.StatusBadRequest, "Invalid invitation ID", &errMsg))
		return
	}

	err = c.service.(*services.SpaceInvitationService).RejectInvitation(uint(invitationId), userId)
	if err != nil {
		errMsg := err.Error()
		ctx.JSON(http.StatusInternalServerError, models.NewErrorResponse(http.StatusInternalServerError, "Failed to reject invitation", &errMsg))
		return
	}

	ctx.JSON(http.StatusOK, models.NewSuccessResponse(
		http.StatusOK,
		"Invitation rejected successfully",
		gin.H{"ok": "yes"},
	))
}

package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/BlenDMinh/dutgrad-server/configs"
	"github.com/BlenDMinh/dutgrad-server/databases/entities"
	"github.com/BlenDMinh/dutgrad-server/databases/repositories"
	"github.com/BlenDMinh/dutgrad-server/helpers"
	"github.com/BlenDMinh/dutgrad-server/models/dtos"
	"github.com/BlenDMinh/dutgrad-server/services"
	"github.com/gin-gonic/gin"
)

type SpaceController struct {
	CrudController[entities.Space, uint]
	service services.SpaceService
}

func NewSpaceController(
	service services.SpaceService,
) *SpaceController {
	crudController := NewCrudController(service)
	return &SpaceController{
		CrudController: *crudController,
		service:        service,
	}
}

func (c *SpaceController) GetPublicSpaces(ctx *gin.Context) {
	params := helpers.GetPaginationParams(ctx, repositories.DefaultPageSize)
	result, err := c.service.GetPublicSpaces(params.Page, params.PageSize)
	if err != nil {
		HandleError(ctx, http.StatusInternalServerError, "Failed to fetch public spaces", err)
		return
	}

	HandleSuccess(ctx, "Public spaces retrieved successfully", gin.H{
		"public_spaces": result.Data,
		"pagination": gin.H{
			"current_page": result.Page,
			"page_size":    result.PageSize,
			"total_pages":  result.TotalPages,
			"total_items":  result.TotalItems,
			"has_next":     result.HasNext,
			"has_prev":     result.HasPrev,
		},
	})
}

func (c *SpaceController) CreateSpace(ctx *gin.Context) {
	model := c.getModel()
	if !HandleBindJSON(ctx, model) {
		return
	}

	userID, ok := ExtractID(ctx, "user_id")
	if !ok {
		return
	}

	createdSpace, err := c.service.CreateSpace(model, userID)
	if err != nil {
		statusCode := http.StatusInternalServerError

		if strings.Contains(err.Error(), "space limit reached") {
			statusCode = http.StatusTooManyRequests
		}

		HandleError(ctx, statusCode, "Failed to create space", err)
		return
	}

	HandleCreated(ctx, "Space created successfully", createdSpace)
}

func (c *SpaceController) GetMembers(ctx *gin.Context) {
	spaceId, ok := ExtractID(ctx, "id")
	if !ok {
		return
	}

	members, err := c.service.GetMembers(spaceId)
	if err != nil {
		HandleError(ctx, http.StatusInternalServerError, "Failed to retrieve members", err)
		return
	}

	HandleSuccess(ctx, "Members retrieved successfully", gin.H{"members": members})
}

func (c *SpaceController) CountSpaceMembers(ctx *gin.Context) {
	spaceId, ok := ExtractID(ctx, "id")
	if !ok {
		return
	}

	count, err := c.service.CountSpaceMembers(spaceId)
	if err != nil {
		HandleError(ctx, http.StatusInternalServerError, "Failed to count members", err)
		return
	}

	HandleSuccess(ctx, "Member count retrieved successfully", gin.H{"count": count})
}

func (c *SpaceController) GetInvitations(ctx *gin.Context) {
	spaceId, ok := ExtractID(ctx, "id")
	if !ok {
		return
	}

	invitations, err := c.service.GetInvitations(spaceId)
	if err != nil {
		HandleError(ctx, http.StatusInternalServerError, "Failed to retrieve invitations", err)
		return
	}

	HandleSuccess(ctx, "Invitations retrieved successfully", gin.H{"invitations": invitations})
}

func (c *SpaceController) GetInvitationLink(ctx *gin.Context) {
	_, ok := ExtractID(ctx, "user_id")
	if !ok {
		return
	}

	spaceId, ok := ExtractID(ctx, "id")
	if !ok {
		return
	}

	var req dtos.GetInvitationLinkRequest
	if !HandleBindJSON(ctx, &req) {
		return
	}

	spaceRoleID := req.SpaceRoleID

	service := c.service
	invitationLink, err := service.GetOrCreateSpaceInvitationLink(spaceId, spaceRoleID)
	if err != nil {
		HandleError(ctx, http.StatusInternalServerError, "Failed to create invitation link", err)
		return
	}

	link, _, err := helpers.GenerateTokenForPayload(
		gin.H{
			"space_id":      invitationLink.SpaceID,
			"space_role_id": invitationLink.SpaceRoleID,
		},
		nil,
	)

	if err != nil {
		HandleError(ctx, http.StatusInternalServerError, "Failed to generate token", err)
		return
	}

	config := configs.GetEnv()
	link = config.WebClientURL + "/invitation?token=" + link

	HandleSuccess(ctx, "Invitation link created successfully", gin.H{"invitation_link": link})
}

func (c *SpaceController) JoinSpace(ctx *gin.Context) {
	token := ctx.Query("token")
	if token == "" {
		HandleError(ctx, http.StatusBadRequest, "Token is required", nil)
		return
	}

	userID, ok := ExtractID(ctx, "user_id")
	if !ok {
		return
	}

	spaceId, err := c.service.JoinSpaceWithToken(token, userID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		message := "Failed to join space"

		if err.Error() == "invalid token" {
			statusCode = http.StatusUnauthorized
			message = "Invalid token"
		} else if err.Error() == "user is already a member of this space" {
			statusCode = http.StatusConflict
			message = "You are already a member of this space"
		}

		HandleError(ctx, statusCode, message, err)
		return
	}

	HandleSuccess(ctx, "Successfully joined the space", gin.H{
		"space_id": spaceId,
	})
}

func (c *SpaceController) GetUserRole(ctx *gin.Context) {
	userID, ok := ExtractID(ctx, "user_id")
	if !ok {
		return
	}

	spaceId, ok := ExtractID(ctx, "id")
	if !ok {
		return
	}

	role, err := c.service.GetUserRole(userID, spaceId)
	if err != nil {
		HandleError(ctx, http.StatusForbidden, "User does not have access to this space", err)
		return
	}

	HandleSuccess(ctx, "User role retrieved successfully", gin.H{"role": role})
}

func (c *SpaceController) InviteUserToSpace(ctx *gin.Context) {
	spaceId, ok := ExtractID(ctx, "id")
	if !ok {
		return
	}

	var req dtos.SpaceInvitationRequest
	if !HandleBindJSON(ctx, &req) {
		return
	}

	inviterId, ok := ExtractID(ctx, "user_id")
	if !ok {
		return
	}

	invitation := entities.SpaceInvitation{
		SpaceID:     spaceId,
		SpaceRoleID: req.SpaceRoleID,
		InviterID:   inviterId,
		Status:      "pending",
		Message:     req.Message,
	}

	if req.InvitedUserID != nil {
		invitation.InvitedUserID = *req.InvitedUserID
	} else {
		userService := services.NewUserService()
		user, err := userService.GetUserByEmail(req.InvitedUserEmail)
		if err != nil {
			HandleError(ctx, http.StatusInternalServerError, "Failed to find user by email", err)
			return
		}
		invitation.InvitedUserID = user.ID
	}

	isMember, err := c.service.IsMemberOfSpace(invitation.InvitedUserID, spaceId)
	if err != nil {
		HandleError(ctx, http.StatusInternalServerError, "Failed to check membership", err)
		return
	}

	if isMember {
		HandleError(ctx, http.StatusBadRequest, "User is already a member of this space", nil)
		return
	}

	_, err = c.service.CreateInvitation(&invitation)
	if err != nil {
		statusCode := http.StatusInternalServerError
		message := "Failed to create invitation"

		if strings.Contains(err.Error(), "uniq_user_space_role") {
			statusCode = http.StatusBadRequest
			message = "User is already a member of this space"
		}

		HandleError(ctx, statusCode, message, err)
		return
	}

	HandleSuccess(ctx, "Invitation sent successfully", gin.H{})
}

func (c *SpaceController) GetSpaceRoles(ctx *gin.Context) {
	roles, err := c.service.GetSpaceRoles()
	if err != nil {
		HandleError(ctx, http.StatusInternalServerError, "Failed to fetch roles", err)
		return
	}

	HandleSuccess(ctx, "Roles retrieved successfully", gin.H{"roles": roles})
}

func (c *SpaceController) JoinPublicSpace(ctx *gin.Context) {
	userID, ok := ExtractID(ctx, "user_id")
	if !ok {
		return
	}

	spaceID, ok := ExtractID(ctx, "id")
	if !ok {
		return
	}

	err := c.service.JoinPublicSpace(spaceID, userID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		message := "Failed to join public space"

		switch err.Error() {
		case "space not found":
			statusCode = http.StatusNotFound
			message = "Space not found"
		case "space is not public":
			statusCode = http.StatusForbidden
			message = "Space is not public"
		case "user is already a member of this space":
			statusCode = http.StatusConflict
			message = "You are already a member of this space"
		}

		HandleError(ctx, statusCode, message, err)
		return
	}

	HandleSuccess(ctx, "Successfully joined the public space", gin.H{})
}

func (c *SpaceController) CountMySpaces(ctx *gin.Context) {
	userID, ok := ExtractID(ctx, "user_id")
	if !ok {
		return
	}

	count, err := c.service.CountSpacesByUserID(userID)
	if err != nil {
		HandleError(ctx, http.StatusInternalServerError, "Failed to count spaces", err)
		return
	}

	ctx.Header("X-Space-Count", fmt.Sprintf("%d", count))
	ctx.Status(http.StatusOK)
}

func (c *SpaceController) GetPopularSpaces(ctx *gin.Context) {
	order := ctx.DefaultQuery("order", "user_count")

	if order != "user_count" {
		HandleError(ctx, http.StatusBadRequest, "Invalid order parameter. Only 'user_count' is supported.", nil)
		return
	}

	popularSpaces, err := c.service.GetPopularSpaces(order)
	if err != nil {
		HandleError(ctx, http.StatusInternalServerError, "Failed to get popular spaces", err)
		return
	}

	HandleSuccess(ctx, "Popular spaces retrieved successfully", gin.H{"popular_spaces": popularSpaces})
}

func (c *SpaceController) Chat(ctx *gin.Context) {
	spaceId, ok := ExtractID(ctx, "id")
	if !ok {
		return
	}

	if c.service.IsAPIRateLimited(spaceId) {
		HandleError(ctx, http.StatusTooManyRequests, "API call limit exceeded for this space", nil)
		return
	}

	var req dtos.ApiChatRequest
	if !HandleBindJSON(ctx, &req) {
		return
	}

	sessionService := services.NewUserQuerySessionService()

	session, err := sessionService.GetById(req.QuerySessionID)
	if err != nil {
		session, err = sessionService.Create(&entities.UserQuerySession{
			SpaceID: spaceId,
		})
		if err != nil {
			HandleError(ctx, http.StatusInternalServerError, "Failed to create session", err)
			return
		}
	}

	ragService := services.NewRAGServerService()
	answer, err := ragService.Chat(session.ID, session.SpaceID, req.Query)
	if err != nil {
		HandleError(ctx, http.StatusInternalServerError, "Failed to get answer", err)
		return
	}

	HandleSuccess(ctx, "Answer retrieved successfully", gin.H{
		"session_id": session.ID,
		"query":      req.Query,
		"answer":     answer,
	})
}

func (c *SpaceController) UpdateUserRole(ctx *gin.Context) {
	spaceId, ok := ExtractID(ctx, "id")
	if !ok {
		return
	}

	memberIdParam := ctx.Param("memberId")
	memberId, err := strconv.ParseUint(memberIdParam, 10, 32)
	if err != nil {
		HandleError(ctx, http.StatusBadRequest, "Invalid member ID", err)
		return
	}

	var req dtos.UpdateRoleRequest
	if !HandleBindJSON(ctx, &req) {
		return
	}

	userID, ok := ExtractID(ctx, "user_id")
	if !ok {
		return
	}

	service := c.service
	err = service.UpdateMemberRole(spaceId, uint(memberId), req.RoleID, userID)
	if err != nil {
		HandleError(ctx, http.StatusInternalServerError, "Failed to update user role", err)
		return
	}

	HandleSuccess(ctx, "User role updated successfully", gin.H{})
}

func (c *SpaceController) RemoveMember(ctx *gin.Context) {
	spaceId, ok := ExtractID(ctx, "id")
	if !ok {
		return
	}

	memberIdParam := ctx.Param("memberId")
	memberId, err := strconv.ParseUint(memberIdParam, 10, 32)
	if err != nil {
		HandleError(ctx, http.StatusBadRequest, "Invalid member ID", err)
		return
	}

	userID, ok := ExtractID(ctx, "user_id")
	if !ok {
		return
	}

	service := c.service
	err = service.RemoveMember(spaceId, uint(memberId), userID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		message := err.Error()

		if strings.Contains(message, "only space owners can remove members") ||
			strings.Contains(message, "cannot remove a space owner") ||
			strings.Contains(message, "you cannot remove yourself") {
			statusCode = http.StatusForbidden
		} else if strings.Contains(message, "not a member of this space") {
			statusCode = http.StatusNotFound
		}

		HandleError(ctx, statusCode, message, nil)
		return
	}

	HandleSuccess(ctx, "Member removed successfully", gin.H{})
}

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
	"github.com/BlenDMinh/dutgrad-server/models"
	"github.com/BlenDMinh/dutgrad-server/models/dtos"
	"github.com/BlenDMinh/dutgrad-server/services"
	"github.com/gin-gonic/gin"
)

type SpaceController struct {
	CrudController[entities.Space, uint]
}

func NewSpaceController() *SpaceController {
	invitationLinkRepo := repositories.NewSpaceInvitationLinkRepository()
	return &SpaceController{
		CrudController: *NewCrudController(services.NewSpaceService(*invitationLinkRepo)),
	}
}

func (c *SpaceController) GetPublicSpaces(ctx *gin.Context) {
	params := helpers.GetPaginationParams(ctx, repositories.DefaultPageSize)
	result, err := c.service.(*services.SpaceService).GetPublicSpaces(params.Page, params.PageSize)
	if err != nil {
		errMsg := err.Error()
		ctx.JSON(
			http.StatusInternalServerError,
			models.NewErrorResponse(
				http.StatusInternalServerError,
				"Failed to fetch public spaces",
				&errMsg,
			),
		)
		return
	}

	ctx.JSON(http.StatusOK, models.NewSuccessResponse(
		http.StatusOK,
		"Success",
		gin.H{
			"public_spaces": result.Data,
			"pagination": gin.H{
				"current_page": result.Page,
				"page_size":    result.PageSize,
				"total_pages":  result.TotalPages,
				"total_items":  result.TotalItems,
				"has_next":     result.HasNext,
				"has_prev":     result.HasPrev,
			},
		},
	))
}

func (c *SpaceController) CreateSpace(ctx *gin.Context) {
	model := c.getModel()
	if err := ctx.ShouldBindJSON(model); err != nil {
		errMsg := err.Error()
		ctx.JSON(400, models.NewErrorResponse(400, "Bad Request", &errMsg))
		return
	}

	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusInternalServerError, models.NewErrorResponse(http.StatusInternalServerError, "User ID not found in context", nil))
		return
	}

	createdSpace, err := c.service.(*services.SpaceService).CreateSpace(model, userID.(uint))
	if err != nil {
		errMsg := err.Error()
		statusCode := http.StatusInternalServerError

		if strings.Contains(errMsg, "space limit reached") {
			statusCode = http.StatusTooManyRequests
		}

		ctx.JSON(statusCode, models.NewErrorResponse(statusCode, "Failed to create space", &errMsg))
		return
	}

	ctx.JSON(201, models.NewSuccessResponse(201, "Created", createdSpace))
}

func (c *SpaceController) GetMembers(ctx *gin.Context) {
	spaceIdParam := ctx.Param("id")
	spaceId, err := strconv.ParseUint(spaceIdParam, 10, 32)

	if err != nil {
		errMsg := err.Error()
		ctx.JSON(
			http.StatusInternalServerError,
			models.NewErrorResponse(
				http.StatusInternalServerError,
				"invalid space id",
				&errMsg,
			),
		)
		return
	}

	members, err := c.service.(*services.SpaceService).GetMembers(uint(spaceId))
	if err != nil {
		errMsg := err.Error()
		ctx.JSON(
			http.StatusInternalServerError,
			models.NewErrorResponse(
				http.StatusInternalServerError,
				"error",
				&errMsg,
			),
		)
		return
	}

	ctx.JSON(http.StatusOK, models.NewSuccessResponse(
		http.StatusOK,
		"Success",
		gin.H{"members": members},
	))
}

func (c *SpaceController) GetInvitations(ctx *gin.Context) {
	spaceIdParam := ctx.Param("id")
	spaceId, err := strconv.ParseUint(spaceIdParam, 10, 32)

	if err != nil {
		errMsg := err.Error()
		ctx.JSON(
			http.StatusInternalServerError,
			models.NewErrorResponse(
				http.StatusInternalServerError,
				"invalid space id",
				&errMsg,
			),
		)
		return
	}

	invitations, err := c.service.(*services.SpaceService).GetInvitations(uint(spaceId))
	if err != nil {
		errMsg := err.Error()
		ctx.JSON(
			http.StatusInternalServerError,
			models.NewErrorResponse(
				http.StatusInternalServerError,
				"error",
				&errMsg,
			),
		)
		return
	}

	ctx.JSON(http.StatusOK, models.NewSuccessResponse(
		http.StatusOK,
		"Success",
		gin.H{"invitations": invitations},
	))
}

func (c *SpaceController) GetInvitationLink(ctx *gin.Context) {
	_, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusInternalServerError, models.NewErrorResponse(http.StatusInternalServerError, "User ID not found in context", nil))
		return
	}

	spaceIdParam := ctx.Param("id")
	spaceId, err := strconv.ParseUint(spaceIdParam, 10, 32)

	var req dtos.GetInvitationLinkRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		errMsg := err.Error()
		ctx.JSON(400, models.NewErrorResponse(400, "Bad Request", &errMsg))
		return
	}

	spaceRoleID := req.SpaceRoleID

	if err != nil {
		errMsg := err.Error()
		ctx.JSON(
			http.StatusInternalServerError,
			models.NewErrorResponse(
				http.StatusInternalServerError,
				"invalid space id",
				&errMsg,
			),
		)
		return
	}

	service := c.service.(*services.SpaceService)
	invitationLink, err := service.GetOrCreateSpaceInvitationLink(uint(spaceId), spaceRoleID)
	if err != nil {
		errMsg := err.Error()
		ctx.JSON(500, models.NewErrorResponse(500, "Internal Server Error", &errMsg))
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
		errMsg := err.Error()
		ctx.JSON(500, models.NewErrorResponse(500, "Internal Server Error", &errMsg))
		return
	}

	config := configs.GetEnv()

	link = config.WebClientURL + "/invitation?token=" + link

	ctx.JSON(http.StatusOK, models.NewSuccessResponse(
		http.StatusOK,
		"Success",
		gin.H{"invitation_link": link},
	))
}

func (c *SpaceController) JoinSpace(ctx *gin.Context) {
	token := ctx.Query("token")
	if token == "" {
		ctx.JSON(http.StatusBadRequest, models.NewErrorResponse(http.StatusBadRequest, "Token is required", nil))
		return
	}

	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusInternalServerError, models.NewErrorResponse(http.StatusInternalServerError, "User ID not found in context", nil))
		return
	}

	spaceId, err := c.service.(*services.SpaceService).JoinSpaceWithToken(token, userID.(uint))
	if err != nil {
		errMsg := err.Error()

		if errMsg == "invalid token" {
			ctx.JSON(http.StatusUnauthorized, models.NewErrorResponse(http.StatusUnauthorized, "Invalid token", &errMsg))
		} else if errMsg == "user is already a member of this space" {
			ctx.JSON(http.StatusConflict, models.NewErrorResponse(http.StatusConflict, "You are already a member of this space", nil))
		} else {
			ctx.JSON(http.StatusInternalServerError, models.NewErrorResponse(http.StatusInternalServerError, "Failed to join space", &errMsg))
		}
		return
	}

	ctx.JSON(http.StatusOK, models.NewSuccessResponse(
		http.StatusOK,
		"Successfully joined the space",
		gin.H{
			"space_id": spaceId,
		},
	))
}

func (c *SpaceController) GetUserRole(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusInternalServerError, models.NewErrorResponse(
			http.StatusInternalServerError,
			"User ID not found in context",
			nil,
		))
		return
	}

	spaceIdParam := ctx.Param("id")
	spaceId, err := strconv.ParseUint(spaceIdParam, 10, 32)
	if err != nil {
		errMsg := err.Error()
		ctx.JSON(
			http.StatusBadRequest,
			models.NewErrorResponse(
				http.StatusBadRequest,
				"Invalid space ID",
				&errMsg,
			),
		)
		return
	}

	role, err := c.service.(*services.SpaceService).GetUserRole(userID.(uint), uint(spaceId))
	if err != nil {
		errMsg := err.Error()
		ctx.JSON(
			http.StatusForbidden,
			models.NewErrorResponse(
				http.StatusForbidden,
				"User does not have access to this space",
				&errMsg,
			),
		)
		return
	}

	ctx.JSON(http.StatusOK, models.NewSuccessResponse(
		http.StatusOK,
		"Success",
		gin.H{"role": role},
	))
}

func (c *SpaceController) InviteUserToSpace(ctx *gin.Context) {
	spaceIdParam := ctx.Param("id")
	spaceId, err := strconv.ParseUint(spaceIdParam, 10, 64)

	if err != nil {
		errMsg := err.Error()
		ctx.JSON(
			http.StatusBadRequest,
			models.NewErrorResponse(
				http.StatusBadRequest,
				"Invalid space id",
				&errMsg,
			),
		)
		return
	}
	var req struct {
		InvitedUserID    *uint  `json:"invited_user_id"`
		InvitedUserEmail string `json:"invited_user_email"`
		SpaceRoleID      uint   `json:"space_role_id" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		errMsg := err.Error()
		ctx.JSON(
			http.StatusBadRequest,
			models.NewErrorResponse(
				http.StatusBadRequest,
				"Invalid request body",
				&errMsg,
			),
		)
		return
	}

	inviterId, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusInternalServerError, models.NewErrorResponse(http.StatusInternalServerError, "User ID not found in context", nil))
		return
	}

	invitation := entities.SpaceInvitation{
		SpaceID:     uint(spaceId),
		SpaceRoleID: req.SpaceRoleID,
		InviterID:   inviterId.(uint),
		Status:      "pending",
	}
	if req.InvitedUserID != nil {
		invitation.InvitedUserID = *req.InvitedUserID
	} else {
		userService := services.NewUserService()
		user, err := userService.GetUserByEmail(req.InvitedUserEmail)
		if err != nil {
			errMsg := err.Error()
			ctx.JSON(
				http.StatusInternalServerError,
				models.NewErrorResponse(
					http.StatusInternalServerError,
					"error",
					&errMsg,
				),
			)
			return
		}
		invitation.InvitedUserID = user.ID
	}

	isMember, err := c.service.(*services.SpaceService).IsMemberOfSpace(invitation.InvitedUserID, uint(spaceId))
	if err != nil {
		errMsg := err.Error()
		ctx.JSON(
			http.StatusInternalServerError,
			models.NewErrorResponse(
				http.StatusInternalServerError,
				"Failed to check membership",
				&errMsg,
			),
		)
		return
	}

	if isMember {
		ctx.JSON(
			http.StatusBadRequest,
			models.NewErrorResponse(
				http.StatusBadRequest,
				"User is already a member of this space",
				nil,
			),
		)
		return
	}

	_, err = c.service.(*services.SpaceService).CreateInvitation(&invitation)
	if err != nil {
		errMsg := err.Error()
		if strings.Contains(errMsg, "uniq_user_space_role") {
			ctx.JSON(
				http.StatusBadRequest,
				models.NewErrorResponse(
					http.StatusBadRequest,
					"User is already a member of this space",
					&errMsg,
				),
			)
			return
		}
		ctx.JSON(
			http.StatusInternalServerError,
			models.NewErrorResponse(
				http.StatusInternalServerError,
				"Unexpected server error",
				&errMsg,
			),
		)
		return
	}

	ctx.JSON(http.StatusOK, models.NewSuccessResponse(
		http.StatusOK,
		"Success",
		gin.H{},
	))
}

func (c *SpaceController) GetSpaceRoles(ctx *gin.Context) {
	roles, err := c.service.(*services.SpaceService).GetSpaceRoles()
	if err != nil {
		errMsg := err.Error()
		ctx.JSON(http.StatusInternalServerError, models.NewErrorResponse(
			http.StatusInternalServerError, "Failed to fetch roles", &errMsg))
		return
	}

	ctx.JSON(http.StatusOK, models.NewSuccessResponse(
		http.StatusOK,
		"Success",
		gin.H{"roles": roles},
	))
}

func (c *SpaceController) JoinPublicSpace(ctx *gin.Context) {
	userIDValue, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, models.NewErrorResponse(http.StatusUnauthorized, "User ID not found in context", nil))
		return
	}
	userID := userIDValue.(uint)

	spaceIdParam := ctx.Param("id")
	spaceID, err := strconv.ParseUint(spaceIdParam, 10, 32)
	if err != nil {
		errMsg := err.Error()
		ctx.JSON(http.StatusBadRequest, models.NewErrorResponse(http.StatusBadRequest, "Invalid space ID", &errMsg))
		return
	}

	err = c.service.(*services.SpaceService).JoinPublicSpace(uint(spaceID), userID)
	if err != nil {
		errMsg := err.Error()
		statusCode := http.StatusInternalServerError
		switch errMsg {
		case "space not found":
			statusCode = http.StatusNotFound
		case "space is not public":
			statusCode = http.StatusForbidden
		case "user is already a member of this space":
			statusCode = http.StatusConflict
		}
		ctx.JSON(statusCode, models.NewErrorResponse(statusCode, errMsg, nil))
		return
	}

	ctx.JSON(http.StatusOK, models.NewSuccessResponse(
		http.StatusOK,
		"Successfully joined the public space",
		gin.H{},
	))
}

func (c *SpaceController) CountMySpaces(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	count, err := c.service.(*services.SpaceService).CountSpacesByUserID(userID.(uint))
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.Header("X-Space-Count", fmt.Sprintf("%d", count))
	ctx.Status(http.StatusOK)
}

func (c *SpaceController) GetPopularSpaces(ctx *gin.Context) {
	order := ctx.DefaultQuery("order", "member_count")

	if order != "member_count" {
		ctx.JSON(http.StatusBadRequest, models.NewErrorResponse(
			http.StatusBadRequest,
			"Invalid order parameter. Only 'member_count' is supported.",
			nil,
		))
		return
	}

	popular_spaces, err := c.service.(*services.SpaceService).GetPopularSpaces(order)

	if err != nil {
		errMsg := err.Error()
		ctx.JSON(http.StatusInternalServerError, models.NewErrorResponse(
			http.StatusInternalServerError,
			"Failed to get popular spaces",
			&errMsg,
		))
		return
	}

	ctx.JSON(http.StatusOK, models.NewSuccessResponse(
		http.StatusOK,
		"Popular spaces retrieved successfully",
		gin.H{"popular_spaces": popular_spaces},
	))
}

func (c *SpaceController) Chat(ctx *gin.Context) {
	spaceIdParam := ctx.Param("id")
	spaceId, err := strconv.ParseUint(spaceIdParam, 10, 32)

	if err != nil {
		errMsg := err.Error()
		ctx.JSON(
			http.StatusInternalServerError,
			models.NewErrorResponse(
				http.StatusInternalServerError,
				"invalid space id",
				&errMsg,
			),
		)
		return
	}
	var req dtos.ApiChatRequest
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

	session, err := sessionService.GetById(req.QuerySessionID)
	if err != nil {
		session, err = sessionService.Create(&entities.UserQuerySession{
			SpaceID: uint(spaceId),
		})
		if err != nil {
			errMsg := err.Error()
			ctx.JSON(
				http.StatusInternalServerError,
				models.NewErrorResponse(
					http.StatusInternalServerError,
					"error",
					&errMsg,
				),
			)
			return
		}
	}

	ragService := services.NewRAGServerService()
	answer, err := ragService.Chat(session.ID, session.SpaceID, req.Query)
	if err != nil {
		errMsg := err.Error()
		ctx.JSON(http.StatusInternalServerError, models.NewErrorResponse(
			http.StatusInternalServerError,
			"Failed to get answer",
			&errMsg,
		))
		return
	}

	ctx.JSON(http.StatusOK, models.NewSuccessResponse(
		http.StatusOK,
		"Answer retrieved successfully",
		&gin.H{
			"session_id": session.ID,
			"query":      req.Query,
			"answer":     answer,
		},
	))
}

func (c *SpaceController) UpdateUserRole(ctx *gin.Context) {
	spaceIdParam := ctx.Param("id")
	spaceId, err := strconv.ParseUint(spaceIdParam, 10, 32)
	if err != nil {
		errMsg := err.Error()
		ctx.JSON(http.StatusBadRequest, models.NewErrorResponse(
			http.StatusBadRequest,
			"Invalid space ID",
			&errMsg,
		))
		return
	}

	memberIdParam := ctx.Param("memberId")
	memberId, err := strconv.ParseUint(memberIdParam, 10, 32)
	if err != nil {
		errMsg := err.Error()
		ctx.JSON(http.StatusBadRequest, models.NewErrorResponse(
			http.StatusBadRequest,
			"Invalid member ID",
			&errMsg,
		))
		return
	}

	var req struct {
		RoleID uint `json:"role_id" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		errMsg := err.Error()
		ctx.JSON(http.StatusBadRequest, models.NewErrorResponse(
			http.StatusBadRequest,
			"Invalid request body",
			&errMsg,
		))
		return
	}

	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusInternalServerError, models.NewErrorResponse(
			http.StatusInternalServerError,
			"User ID not found in context",
			nil,
		))
		return
	}

	service := c.service.(*services.SpaceService)
	err = service.UpdateMemberRole(uint(spaceId), uint(memberId), req.RoleID, userID.(uint))
	if err != nil {
		errMsg := err.Error()
		ctx.JSON(http.StatusInternalServerError, models.NewErrorResponse(
			http.StatusInternalServerError,
			"Failed to update user role",
			&errMsg,
		))
		return
	}

	ctx.JSON(http.StatusOK, models.NewSuccessResponse(
		http.StatusOK,
		"User role updated successfully",
		gin.H{},
	))
}

func (c *SpaceController) RemoveMember(ctx *gin.Context) {
	spaceIdParam := ctx.Param("id")
	spaceId, err := strconv.ParseUint(spaceIdParam, 10, 32)
	if err != nil {
		errMsg := err.Error()
		ctx.JSON(http.StatusBadRequest, models.NewErrorResponse(
			http.StatusBadRequest,
			"Invalid space ID",
			&errMsg,
		))
		return
	}

	memberIdParam := ctx.Param("memberId")
	memberId, err := strconv.ParseUint(memberIdParam, 10, 32)
	if err != nil {
		errMsg := err.Error()
		ctx.JSON(http.StatusBadRequest, models.NewErrorResponse(
			http.StatusBadRequest,
			"Invalid member ID",
			&errMsg,
		))
		return
	}

	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusInternalServerError, models.NewErrorResponse(
			http.StatusInternalServerError,
			"User ID not found in context",
			nil,
		))
		return
	}

	service := c.service.(*services.SpaceService)
	err = service.RemoveMember(uint(spaceId), uint(memberId), userID.(uint))
	if err != nil {
		errMsg := err.Error()
		statusCode := http.StatusInternalServerError

		if strings.Contains(errMsg, "only space owners can remove members") ||
			strings.Contains(errMsg, "cannot remove a space owner") ||
			strings.Contains(errMsg, "you cannot remove yourself") {
			statusCode = http.StatusForbidden
		} else if strings.Contains(errMsg, "not a member of this space") {
			statusCode = http.StatusNotFound
		}

		ctx.JSON(statusCode, models.NewErrorResponse(
			statusCode,
			errMsg,
			nil,
		))
		return
	}

	ctx.JSON(http.StatusOK, models.NewSuccessResponse(
		http.StatusOK,
		"Member removed successfully",
		gin.H{},
	))
}

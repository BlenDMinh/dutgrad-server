package controllers

import (
	"net/http"

	"github.com/BlenDMinh/dutgrad-server/databases/entities"
	"github.com/BlenDMinh/dutgrad-server/models/dtos"
	"github.com/BlenDMinh/dutgrad-server/services"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	CrudController[entities.User, uint]
	service services.UserService
}

func NewUserController(
	service services.UserService,
) *UserController {
	crudController := NewCrudController(service)
	return &UserController{
		CrudController: *crudController,
		service:        service,
	}
}

func (uc *UserController) GetCurrentUser(ctx *gin.Context) {
	userID, ok := ExtractID(ctx, "user_id")
	if !ok {
		return
	}

	user, err := uc.service.GetById(userID)
	if err != nil {
		HandleError(ctx, http.StatusInternalServerError, "Failed to retrieve user", err)
		return
	}

	HandleSuccess(ctx, "User retrieved successfully", user)
}

func (c *UserController) GetMySpaces(ctx *gin.Context) {
	userId, ok := ExtractID(ctx, "user_id")
	if !ok {
		return
	}

	spaces, err := c.service.GetSpacesByUserId(userId)
	if err != nil {
		HandleError(ctx, http.StatusInternalServerError, "Failed to retrieve spaces", err)
		return
	}

	HandleSuccess(ctx, "Success", dtos.UserSpaceListResponse{
		Spaces: spaces,
	})
}

func (c *UserController) GetUserSpaces(ctx *gin.Context) {
	userId, ok := ExtractID(ctx, "id")
	if !ok {
		return
	}

	spaces, err := c.service.GetSpacesByUserId(userId)
	if err != nil {
		HandleError(ctx, http.StatusInternalServerError, "Failed to retrieve spaces", err)
		return
	}

	HandleSuccess(ctx, "Success", dtos.UserSpaceListResponse{
		Spaces: spaces,
	})
}

func (c *UserController) GetMyInvitations(ctx *gin.Context) {
	userId, ok := ExtractID(ctx, "user_id")
	if !ok {
		return
	}

	invitations, err := c.service.GetInvitationsByUserId(userId)
	if err != nil {
		HandleError(ctx, http.StatusInternalServerError, "Failed to retrieve invitations", err)
		return
	}

	HandleSuccess(ctx, "Invitations retrieved successfully", dtos.SpaceInvitationListResponse{
		Invitations: invitations,
	})
}

func (c *UserController) SearchUsers(ctx *gin.Context) {
	query := ctx.Query("query")
	if query == "" {
		HandleError(ctx, http.StatusBadRequest, "Search query cannot be empty", nil)
		return
	}

	users, err := c.service.SearchUsers(query)
	if err != nil {
		HandleError(ctx, http.StatusInternalServerError, "Failed to search users", err)
		return
	}

	HandleSuccess(ctx, "Users retrieved successfully", dtos.UserListResponse{
		Users: users,
	})
}

func (c *UserController) GetUserTier(ctx *gin.Context) {
	userID, ok := ExtractID(ctx, "user_id")
	if !ok {
		return
	}

	tierUsage, err := c.service.GetUserTierUsage(userID)
	if err != nil {
		HandleError(ctx, http.StatusInternalServerError, "Failed to fetch user tier and usage", err)
		return
	}

	HandleSuccess(ctx, "User tier and usage fetched successfully", dtos.UserTierUsageResponse{
		Tier:  tierUsage.Tier,
		Usage: tierUsage.Usage,
	})
}

func (c *UserController) UpdatePassword(ctx *gin.Context) {
	userID, ok := ExtractID(ctx, "user_id")
	if !ok {
		return
	}
	var req struct {
		CurrentPassword string `json:"currentPassword"`
		NewPassword     string `json:"newPassword"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		HandleError(ctx, http.StatusBadRequest, "Invalid request data", err)
		return
	}
	if len(req.NewPassword) < 6 {
		HandleError(ctx, http.StatusBadRequest, "New password must be at least 6 characters long", nil)
		return
	}

	if req.CurrentPassword == req.NewPassword {
		HandleError(ctx, http.StatusBadRequest, "New password must be different from current password", nil)
		return
	}

	authService := services.NewAuthService()
	err := authService.ChangePassword(userID, req.CurrentPassword, req.NewPassword)
	if err != nil {
		HandleError(ctx, http.StatusBadRequest, "Failed to update password", err)
		return
	}

	HandleSuccess(ctx, "Password updated successfully", nil)
}

func (c *UserController) GetUserAuthMethod(ctx *gin.Context) {
	userID, ok := ExtractID(ctx, "user_id")
	if !ok {
		return
	}

	authService := services.NewAuthService()
	authMethods, err := authService.GetUserAuthMethods(userID)
	if err != nil {
		HandleError(ctx, http.StatusInternalServerError, "Failed to retrieve auth methods", err)
		return
	}

	HandleSuccess(ctx, "Auth methods retrieved successfully", authMethods)
}

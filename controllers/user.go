package controllers

import (
	"net/http"
	"strconv"

	"github.com/BlenDMinh/dutgrad-server/databases/entities"
	"github.com/BlenDMinh/dutgrad-server/databases/repositories"
	"github.com/BlenDMinh/dutgrad-server/models"
	"github.com/BlenDMinh/dutgrad-server/services"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	CrudController[entities.User, uint]
}

func NewUserController() *UserController {
	return &UserController{
		CrudController: *NewCrudController(services.NewUserService()),
	}
}

func (uc *UserController) GetCurrentUser(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusInternalServerError, models.NewErrorResponse(http.StatusInternalServerError, "User ID not found in context", nil))
		return
	}

	user, err := uc.service.GetById(userID.(uint))
	if err != nil {
		errMsg := err.Error()
		ctx.JSON(http.StatusInternalServerError, models.NewErrorResponse(http.StatusInternalServerError, "Failed to retrieve user", &errMsg))
		return
	}

	ctx.JSON(http.StatusOK, models.NewSuccessResponse(http.StatusOK, "User retrieved successfully", user))
}

func (c *UserController) GetMySpaces(ctx *gin.Context) {
	userId, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusInternalServerError, models.NewErrorResponse(http.StatusInternalServerError, "User ID not found in context", nil))
		return
	}
	spaces, err := c.service.(*services.UserService).GetSpacesByUserId(userId.(uint))
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
		gin.H{"spaces": spaces},
	))
}

func (c *UserController) GetUserSpaces(ctx *gin.Context) {
	userIdParam := ctx.Param("user_id")
	userId, err := strconv.ParseUint(userIdParam, 10, 32)
	if err != nil {
		errMsg := err.Error()
		ctx.JSON(
			http.StatusInternalServerError,
			models.NewErrorResponse(
				http.StatusInternalServerError,
				"Invalid user ID",
				&errMsg,
			),
		)
		return
	}

	spaces, err := c.service.(*services.UserService).GetSpacesByUserId(uint(userId))
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
		gin.H{"spaces": spaces},
	))
}

func (c *UserController) GetMyInvitations(ctx *gin.Context) {
	userId, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusInternalServerError, models.NewErrorResponse(http.StatusInternalServerError, "User ID not found in context", nil))
		return
	}

	invitations, err := c.service.(*services.UserService).GetInvitationsByUserId(userId.(uint))
	if err != nil {
		errMsg := err.Error()
		ctx.JSON(
			http.StatusInternalServerError,
			models.NewErrorResponse(
				http.StatusInternalServerError,
				"Failed to retrieve invitations",
				&errMsg,
			),
		)
		return
	}

	ctx.JSON(http.StatusOK, models.NewSuccessResponse(
		http.StatusOK,
		"Invitations retrieved successfully",
		gin.H{"invitations": invitations},
	))
}

func (c *UserController) SearchUsers(ctx *gin.Context) {
	query := ctx.Query("query")
	if query == "" {
		ctx.JSON(
			http.StatusBadRequest,
			models.NewErrorResponse(
				http.StatusBadRequest,
				"Search query cannot be empty",
				nil,
			),
		)
		return
	}

	users, err := c.service.(*services.UserService).SearchUsers(query)
	if err != nil {
		errMsg := err.Error()
		ctx.JSON(
			http.StatusInternalServerError,
			models.NewErrorResponse(
				http.StatusInternalServerError,
				"Failed to search users",
				&errMsg,
			),
		)
		return
	}

	ctx.JSON(http.StatusOK, models.NewSuccessResponse(
		http.StatusOK,
		"Users retrieved successfully",
		gin.H{"users": users},
	))
}

func (c *UserController) GetUserTier(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, models.NewErrorResponse(
			http.StatusUnauthorized,
			"User ID not found in context",
			nil,
		))
		return
	}

	service := c.service.(*services.UserService)
	tierUsage, err := service.GetUserTierUsage(userID.(uint))
	if err != nil {
		errMsg := err.Error()
		ctx.JSON(http.StatusInternalServerError, models.NewErrorResponse(
			http.StatusInternalServerError,
			"Failed to fetch user tier and usage",
			&errMsg,
		))
		return
	}

	ctx.JSON(http.StatusOK, models.NewSuccessResponse(
		http.StatusOK,
		"User tier and usage fetched successfully",
		gin.H{
			"tier":  tierUsage.Tier,
			"usage": tierUsage.Usage,
		},
	))
}

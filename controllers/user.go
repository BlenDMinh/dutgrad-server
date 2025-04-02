package controllers

import (
	"net/http"
	"strconv"

	"github.com/BlenDMinh/dutgrad-server/databases/entities"
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
	userID, exists := ctx.Get("userID")
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

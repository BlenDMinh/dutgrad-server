package controllers

import (
	"net/http"

	"github.com/BlenDMinh/dutgrad-server/databases/entities"
	"github.com/BlenDMinh/dutgrad-server/models"
	"github.com/BlenDMinh/dutgrad-server/services"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	CrudController[entities.User]
	userService services.UserService
}

func NewUserController() *UserController {
	return &UserController{
		userService: services.UserService{},
	}
}

func (uc *UserController) GetCurrentUser(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusInternalServerError, models.NewErrorResponse(http.StatusInternalServerError, "User ID not found in context", nil))
		return
	}

	user, err := uc.userService.GetUserByID(userID.(uint))
	if err != nil {
		errMsg := err.Error()
		ctx.JSON(http.StatusInternalServerError, models.NewErrorResponse(http.StatusInternalServerError, "Failed to retrieve user", &errMsg))
		return
	}

	ctx.JSON(http.StatusOK, models.NewSuccessResponse(http.StatusOK, "User retrieved successfully", user))
}

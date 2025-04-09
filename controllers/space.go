package controllers

import (
	"net/http"

	"github.com/BlenDMinh/dutgrad-server/databases"
	"github.com/BlenDMinh/dutgrad-server/databases/entities"
	"github.com/BlenDMinh/dutgrad-server/models"
	"github.com/BlenDMinh/dutgrad-server/services"
	"github.com/gin-gonic/gin"
)

type SpaceController struct {
	CrudController[entities.Space, uint]
}

func NewSpaceController() *SpaceController {
	return &SpaceController{
		CrudController: *NewCrudController(services.NewSpaceService()),
	}
}

func (c *SpaceController) GetPublicSpaces(ctx *gin.Context) {
	spaces, err := c.service.(*services.SpaceService).GetPublicSpaces()
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
		gin.H{"public_spaces": spaces},
	))
}

func (c *SpaceController) CreateSpace(ctx *gin.Context) {
	model := c.getModel()
	if err := ctx.ShouldBindJSON(model); err != nil {
		errMsg := err.Error()
		ctx.JSON(400, models.NewErrorResponse(400, "Bad Request", &errMsg))
		return
	}

	createdSpace, err := c.service.Create(model)
	if err != nil {
		errMsg := err.Error()
		ctx.JSON(500, models.NewErrorResponse(500, "Internal Server Error", &errMsg))
		return
	}

	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusInternalServerError, models.NewErrorResponse(http.StatusInternalServerError, "User ID not found in context", nil))
		return
	}

	spaceRoleID := uint(entities.Owner)

	spaceUser := entities.SpaceUser{
		UserID:      userID.(uint),
		SpaceID:     createdSpace.ID,
		SpaceRoleID: &spaceRoleID,
	}

	db := databases.GetDB()
	if err := db.Create(&spaceUser).Error; err != nil {
		errMsg := err.Error()
		ctx.JSON(500, models.NewErrorResponse(500, "Failed to create SpaceUser", &errMsg))
		return
	}

	ctx.JSON(201, models.NewSuccessResponse(201, "Created", createdSpace))
}

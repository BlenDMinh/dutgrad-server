package controllers

import (
	"net/http"
	"strconv"

	"github.com/BlenDMinh/dutgrad-server/configs"
	"github.com/BlenDMinh/dutgrad-server/databases"
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
	// Validate token parameter
	token := ctx.Query("token")
	if token == "" {
		ctx.JSON(http.StatusBadRequest, models.NewErrorResponse(http.StatusBadRequest, "Token is required", nil))
		return
	}

	// Get user ID from context
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusInternalServerError, models.NewErrorResponse(http.StatusInternalServerError, "User ID not found in context", nil))
		return
	}

	// Call service to handle business logic
	err := c.service.(*services.SpaceService).JoinSpaceWithToken(token, userID.(uint))
	if err != nil {
		errMsg := err.Error()

		// Handle different error types with appropriate status codes
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
		nil,
	))
}

func (c *SpaceController) GetUserRole(ctx *gin.Context) {
	// Get user ID from context
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusInternalServerError, models.NewErrorResponse(
			http.StatusInternalServerError,
			"User ID not found in context",
			nil,
		))
		return
	}

	// Get space ID from path parameter
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

	// Get user role in the space
	role, err := c.service.(*services.SpaceService).GetUserRole(userID.(uint), uint(spaceId))
	if err != nil {
		errMsg := err.Error()
		// Return 403 if user has no role in this space
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

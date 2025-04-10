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

	_, err = c.service.(*services.SpaceService).CreateInvitation(&invitation)
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

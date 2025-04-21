package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/BlenDMinh/dutgrad-server/databases/entities"
	"github.com/BlenDMinh/dutgrad-server/helpers"
	"github.com/BlenDMinh/dutgrad-server/services"
	"github.com/gin-gonic/gin"
)

type SpaceApiKeyController struct {
	CrudController[entities.SpaceAPIKey, uint]
}

func NewSpaceApiKeyController() *SpaceApiKeyController {
	return &SpaceApiKeyController{
		CrudController: *NewCrudController(services.NewSpaceApiKeyService()),
	}
}

func (ctrl *SpaceApiKeyController) Create(c *gin.Context) {
	spaceID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid space id"})
		return
	}

	var input struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	apiKey := entities.SpaceAPIKey{
		Name:        input.Name,
		Description: input.Description,
		SpaceID:     uint(spaceID),
	}

	created, err := ctrl.service.(*services.SpaceApiKeyService).Create(&apiKey)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, created)
}

func (ctrl *SpaceApiKeyController) List(c *gin.Context) {
	spaceId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid space id"})
		return
	}

	items, err := ctrl.service.(*services.SpaceApiKeyService).GetAllBySpaceID(uint(spaceId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, items)
}

func (ctrl *SpaceApiKeyController) GetOne(c *gin.Context) {
	spaceID, err1 := strconv.Atoi(c.Param("id"))
	keyID, err2 := strconv.Atoi(c.Param("keyId"))
	if err1 != nil || err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ids"})
		return
	}

	item, err := ctrl.service.GetById(uint(keyID))
	if err != nil || item == nil || item.SpaceID != uint(spaceID) {
		c.JSON(http.StatusNotFound, gin.H{"error": "API key not found"})
		return
	}

	token, _, err := helpers.GenerateTokenForPayload(map[string]interface{}{
		"space_id": item.SpaceID,
		"key_id":   item.ID,
	}, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate JWT"})
		return
	}

	response := gin.H{
		"id":          item.ID,
		"name":        item.Name,
		"description": item.Description,
		"space_id":    item.SpaceID,
		"token":       token,
	}

	c.JSON(http.StatusOK, response)
}

func (ctrl *SpaceApiKeyController) Delete(c *gin.Context) {
	spaceID, err1 := strconv.Atoi(c.Param("id"))
	keyID, err2 := strconv.Atoi(c.Param("keyId"))
	if err1 != nil || err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ids"})
		return
	}

	item, err := ctrl.service.GetById(uint(keyID))
	if err != nil || item == nil || item.SpaceID != uint(spaceID) {
		c.JSON(http.StatusNotFound, gin.H{"error": "API key not found"})
		return
	}

	err = ctrl.service.Delete(uint(keyID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "API key deleted successfully"})
}

func VerifyBearerToken(tokenString string) (*entities.SpaceAPIKey, error) {
	payload, err := helpers.VerifyTokenForPayload(tokenString)
	if err != nil {
		return nil, err
	}

	keyIDFloat, ok := (*payload)["key_id"].(float64)
	if !ok {
		return nil, errors.New("invalid key_id in token")
	}

	keyID := uint(keyIDFloat)
	service := services.NewSpaceApiKeyService()
	item, err := service.GetById(keyID)
	if err != nil || item == nil {
		return nil, errors.New("invalid token: key not found")
	}

	return item, nil
}

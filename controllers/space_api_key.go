package controllers

import (
	"errors"
	"net/http"

	"github.com/BlenDMinh/dutgrad-server/databases/entities"
	"github.com/BlenDMinh/dutgrad-server/helpers"
	"github.com/BlenDMinh/dutgrad-server/models/dtos"
	"github.com/BlenDMinh/dutgrad-server/services"
	"github.com/gin-gonic/gin"
)

type SpaceApiKeyController struct {
	CrudController[entities.SpaceAPIKey, uint]
	service services.SpaceApiKeyService
}

func NewSpaceApiKeyController(
	service services.SpaceApiKeyService,
) *SpaceApiKeyController {
	crudController := NewCrudController(service)
	return &SpaceApiKeyController{
		CrudController: *crudController,
		service:        service,
	}
}

func (c *SpaceApiKeyController) Create(ctx *gin.Context) {
	spaceID, ok := ExtractID(ctx, "id")
	if !ok {
		return
	}

	var input dtos.CreateApiKeyRequest
	if !HandleBindJSON(ctx, &input) {
		return
	}

	apiKey := entities.SpaceAPIKey{
		Name:        input.Name,
		Description: input.Description,
		SpaceID:     spaceID,
	}

	created, err := c.service.Create(&apiKey)
	if err != nil {
		HandleError(ctx, http.StatusInternalServerError, "Failed to create API key", err)
		return
	}

	token, _, err := helpers.GenerateTokenForPayload(map[string]interface{}{
		"space_id": created.SpaceID,
		"key_id":   created.ID,
	}, nil)
	if err != nil {
		HandleError(ctx, http.StatusInternalServerError, "Failed to generate JWT", err)
		return
	}

	response := dtos.ApiKeyResponse{
		ID:          created.ID,
		Name:        created.Name,
		Description: created.Description,
		SpaceID:     created.SpaceID,
		Token:       token,
	}

	HandleCreated(ctx, "API key created successfully", response)
}

func (c *SpaceApiKeyController) List(ctx *gin.Context) {
	spaceId, ok := ExtractID(ctx, "id")
	if !ok {
		return
	}

	items, err := c.service.GetAllBySpaceID(spaceId)
	if err != nil {
		HandleError(ctx, http.StatusInternalServerError, "Failed to retrieve API keys", err)
		return
	}

	HandleSuccess(ctx, "API keys retrieved successfully", gin.H{"API": items})
}

func (c *SpaceApiKeyController) GetOne(ctx *gin.Context) {
	spaceID, ok1 := ExtractID(ctx, "id")
	if !ok1 {
		return
	}

	keyID, ok2 := ExtractID(ctx, "keyId")
	if !ok2 {
		return
	}

	item, err := c.service.GetById(keyID)
	if err != nil || item == nil || item.SpaceID != spaceID {
		HandleError(ctx, http.StatusNotFound, "API key not found", err)
		return
	}

	token, _, err := helpers.GenerateTokenForPayload(map[string]interface{}{
		"space_id": item.SpaceID,
		"key_id":   item.ID,
	}, nil)
	if err != nil {
		HandleError(ctx, http.StatusInternalServerError, "Failed to generate JWT", err)
		return
	}

	response := dtos.ApiKeyResponse{
		ID:          item.ID,
		Name:        item.Name,
		Description: item.Description,
		SpaceID:     item.SpaceID,
		Token:       token,
	}

	HandleSuccess(ctx, "API key retrieved successfully", gin.H{"API": response})
}

func (c *SpaceApiKeyController) Delete(ctx *gin.Context) {
	spaceID, ok1 := ExtractID(ctx, "id")
	if !ok1 {
		return
	}

	keyID, ok2 := ExtractID(ctx, "keyId")
	if !ok2 {
		return
	}

	item, err := c.service.GetById(keyID)
	if err != nil || item == nil || item.SpaceID != spaceID {
		HandleError(ctx, http.StatusNotFound, "API key not found", err)
		return
	}

	err = c.service.Delete(keyID)
	if err != nil {
		HandleError(ctx, http.StatusInternalServerError, "Failed to delete API key", err)
		return
	}

	HandleSuccess(ctx, "API key deleted successfully", nil)
}

func VerifySpaceAPIKey(tokenString string) (*entities.SpaceAPIKey, error) {
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

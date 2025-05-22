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
}

func NewSpaceApiKeyController() *SpaceApiKeyController {
	return &SpaceApiKeyController{
		CrudController: *NewCrudController(services.NewSpaceApiKeyService()),
	}
}

func (ctrl *SpaceApiKeyController) Create(c *gin.Context) {
	spaceID, ok := ExtractID(c, "id")
	if !ok {
		return
	}

	var input dtos.CreateApiKeyRequest
	if !HandleBindJSON(c, &input) {
		return
	}

	apiKey := entities.SpaceAPIKey{
		Name:        input.Name,
		Description: input.Description,
		SpaceID:     spaceID,
	}

	created, err := ctrl.service.(*services.SpaceApiKeyService).Create(&apiKey)
	if err != nil {
		HandleError(c, http.StatusInternalServerError, "Failed to create API key", err)
		return
	}

	token, _, err := helpers.GenerateTokenForPayload(map[string]interface{}{
		"space_id": created.SpaceID,
		"key_id":   created.ID,
	}, nil)
	if err != nil {
		HandleError(c, http.StatusInternalServerError, "Failed to generate JWT", err)
		return
	}

	response := dtos.ApiKeyResponse{
		ID:          created.ID,
		Name:        created.Name,
		Description: created.Description,
		SpaceID:     created.SpaceID,
		Token:       token,
	}

	HandleCreated(c, "API key created successfully", response)
}

func (ctrl *SpaceApiKeyController) List(c *gin.Context) {
	spaceId, ok := ExtractID(c, "id")
	if !ok {
		return
	}

	items, err := ctrl.service.(*services.SpaceApiKeyService).GetAllBySpaceID(spaceId)
	if err != nil {
		HandleError(c, http.StatusInternalServerError, "Failed to retrieve API keys", err)
		return
	}

	HandleSuccess(c, "API keys retrieved successfully", gin.H{"API": items})
}

func (ctrl *SpaceApiKeyController) GetOne(c *gin.Context) {
	spaceID, ok1 := ExtractID(c, "id")
	if !ok1 {
		return
	}

	keyID, ok2 := ExtractID(c, "keyId")
	if !ok2 {
		return
	}

	item, err := ctrl.service.GetById(keyID)
	if err != nil || item == nil || item.SpaceID != spaceID {
		HandleError(c, http.StatusNotFound, "API key not found", err)
		return
	}

	token, _, err := helpers.GenerateTokenForPayload(map[string]interface{}{
		"space_id": item.SpaceID,
		"key_id":   item.ID,
	}, nil)
	if err != nil {
		HandleError(c, http.StatusInternalServerError, "Failed to generate JWT", err)
		return
	}

	response := dtos.ApiKeyResponse{
		ID:          item.ID,
		Name:        item.Name,
		Description: item.Description,
		SpaceID:     item.SpaceID,
		Token:       token,
	}

	HandleSuccess(c, "API key retrieved successfully", gin.H{"API": response})
}

func (ctrl *SpaceApiKeyController) Delete(c *gin.Context) {
	spaceID, ok1 := ExtractID(c, "id")
	if !ok1 {
		return
	}

	keyID, ok2 := ExtractID(c, "keyId")
	if !ok2 {
		return
	}

	item, err := ctrl.service.GetById(keyID)
	if err != nil || item == nil || item.SpaceID != spaceID {
		HandleError(c, http.StatusNotFound, "API key not found", err)
		return
	}

	err = ctrl.service.Delete(keyID)
	if err != nil {
		HandleError(c, http.StatusInternalServerError, "Failed to delete API key", err)
		return
	}

	HandleSuccess(c, "API key deleted successfully", nil)
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

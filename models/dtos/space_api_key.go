package dtos

type CreateApiKeyRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type ApiKeyResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	SpaceID     uint   `json:"space_id"`
	Token       string `json:"token"`
}

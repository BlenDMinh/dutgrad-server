package dtos

// CreateApiKeyRequest represents the request body for creating a new API key
type CreateApiKeyRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

// ApiKeyResponse represents the response after creating or retrieving an API key
type ApiKeyResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	SpaceID     uint   `json:"space_id"`
	Token       string `json:"token"`
}

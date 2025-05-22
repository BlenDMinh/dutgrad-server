package dtos

import "time"

// DocumentUploadRequest represents the request body for uploading a document
type DocumentUploadRequest struct {
	SpaceID     uint   `form:"space_id" binding:"required"`
	Description string `form:"description"`
}

// DocumentResponse represents the response after creating or retrieving a document
type DocumentResponse struct {
	ID               uint      `json:"id"`
	Name             string    `json:"name"`
	Description      string    `json:"description"`
	SpaceID          uint      `json:"space_id"`
	MimeType         string    `json:"mime_type"`
	URL              string    `json:"url"`
	ProcessingStatus string    `json:"processing_status"`
	SizeKb           int       `json:"size_kb"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

package dtos

type BeginChatSessionRequest struct {
	SpaceID uint `json:"space_id" binding:"required"`
}

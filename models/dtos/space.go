package dtos

type GetInvitationLinkRequest struct {
	SpaceRoleID uint `json:"space_role_id" binding:"required"`
}

type ApiChatRequest struct {
	QuerySessionID uint   `json:"query_session_id"`
	Query          string `json:"query" binding:"required"`
}

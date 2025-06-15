package dtos

import (
	"time"

	"github.com/BlenDMinh/dutgrad-server/databases/entities"
)

type GetInvitationLinkRequest struct {
	SpaceRoleID uint `json:"space_role_id" binding:"required"`
}

type ApiChatRequest struct {
	QuerySessionID uint   `json:"query_session_id"`
	Query          string `json:"query" binding:"required"`
}

type UserSpaceDTO struct {
	ID              uint               `json:"id"`
	Name            string             `json:"name"`
	Description     string             `json:"description"`
	PrivacyStatus   bool               `json:"privacy_status"`
	DocumentLimit   int                `json:"document_limit"`
	FileSizeLimitKb int                `json:"file_size_limit_kb"`
	ApiCallLimit    int                `json:"api_call_limit"`
	CreatedAt       time.Time          `json:"created_at"`
	UpdatedAt       time.Time          `json:"updated_at"`
	Role            entities.SpaceRole `json:"role"`
	UserCount       int                `json:"user_count"`
}

type SpaceInvitationRequest struct {
	InvitedUserID    *uint  `json:"invited_user_id"`
	InvitedUserEmail string `json:"invited_user_email"`
	SpaceRoleID      uint   `json:"space_role_id" binding:"required"`
	Message          string `json:"message"`
}

type UpdateRoleRequest struct {
	RoleID uint `json:"role_id" binding:"required"`
}

type SpaceListResponse struct {
	Spaces []entities.Space `json:"spaces"`
}

type UserSpaceListResponse struct {
	Spaces []UserSpaceDTO `json:"spaces"`
}

type SpaceUsage struct {
	SpaceID                uint  `json:"space_id"`
	ChatAPICallsUsageDaily int64 `json:"chat_api_calls_usage_daily"`
}

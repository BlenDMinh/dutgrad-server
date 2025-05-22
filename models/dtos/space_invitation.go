package dtos

import "github.com/BlenDMinh/dutgrad-server/databases/entities"

// SpaceInvitationResponse represents a space invitation response
type SpaceInvitationResponse struct {
	ID        uint   `json:"id"`
	UserID    uint   `json:"user_id"`
	SpaceID   uint   `json:"space_id"`
	RoleID    uint   `json:"role_id"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type SpaceInvitationListResponse struct {
	Invitations []entities.SpaceInvitation `json:"invitations"`
}

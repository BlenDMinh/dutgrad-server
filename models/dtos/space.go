package dtos

type GetInvitationLinkRequest struct {
	SpaceRoleID uint   `json:"space_role_id" binding:"required"`
}

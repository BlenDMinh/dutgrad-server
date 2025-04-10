package dtos

type GetInvitationLinkRequest struct {
	SpaceID uint   `json:"space_id" binding:"required"`
	SpaceRoleID uint   `json:"space_role_id" binding:"required"`
}

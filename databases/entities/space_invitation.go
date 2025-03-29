package entities

import "time"

type SpaceInvitation struct {
	ID            uint       `json:"id"`
	SpaceID       uint       `json:"space_id"`
	Space         *Space     `json:"space" gorm:"foreignKey:SpaceID"`
	SpaceRoleID   uint       `json:"space_role_id"`
	SpaceRole     *SpaceRole `json:"space_role" gorm:"foreignKey:SpaceRoleID"`
	InvitedUserID uint       `json:"invited_user_id"`
	InvitedUser   *User      `json:"invited_user" gorm:"foreignKey:InvitedUserID"`
	InviterID     uint       `json:"inviter_id"`
	Inviter       *User      `json:"inviter" gorm:"foreignKey:InviterID"`
	Status        string     `json:"status"` // e.g., "pending", "accepted", "declined"
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

func (s SpaceInvitation) GetIdType() string {
	return "uint"
}

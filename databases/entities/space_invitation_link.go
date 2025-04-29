package entities

import "time"

type SpaceInvitationLink struct {
	ID          uint       `json:"id" gorm:"primaryKey"`
	SpaceID     uint       `json:"space_id"`
	Space       *Space     `json:"space" gorm:"foreignKey:SpaceID"`
	SpaceRoleID uint       `json:"space_role_id"`
	SpaceRole   *SpaceRole `json:"space_role" gorm:"foreignKey:SpaceRoleID"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func (s SpaceInvitationLink) GetIdType() string {
	return "uint"
}

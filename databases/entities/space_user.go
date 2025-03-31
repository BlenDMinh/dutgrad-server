package entities

import "time"

type SpaceUser struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	UserID      uint      `json:"user_id" gorm:"not null;index"`
	SpaceID     uint      `json:"space_id" gorm:"not null;index"`
	SpaceRoleID *uint     `json:"space_role_id" gorm:"index"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	User        User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;" json:"user"`
	Space       Space     `gorm:"foreignKey:SpaceID;constraint:OnDelete:CASCADE;" json:"space"`
	SpaceRole   SpaceRole `gorm:"foreignKey:SpaceRoleID;constraint:OnDelete:SET NULL;" json:"space_role"`
}

func (s SpaceUser) GetIdType() string {
	return "uint"
}

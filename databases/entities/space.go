package entities

import "time"

type Space struct {
	ID            uint               `gorm:"primaryKey" json:"id"`
	Name          string             `gorm:"type:varchar(255);not null" json:"name"`
	Description   string             `gorm:"type:text" json:"description"`
	PrivacyStatus bool               `json:"privacy_status"`
	CreatedAt     time.Time          `json:"created_at"`
	UpdatedAt     time.Time          `json:"updated_at"`
	Documents     []Document         `gorm:"foreignKey:SpaceID" json:"documents"`
	Sessions      []UserQuerySession `gorm:"foreignKey:SpaceID" json:"sessions"`
}

func (s Space) GetIdType() string {
	return "uint"
}

package entities

import "time"

type User struct {
	ID         uint               `json:"id" gorm:"primaryKey"`
	Username   string             `json:"username" gorm:"size:100;not null;uniqueIndex"`
	Email      *string            `json:"email" gorm:"size:100;uniqueIndex"`
	Active     bool               `json:"active" gorm:"default:true"`
	MFAEnabled bool               `json:"mfa_enabled" gorm:"default:false"`
	TierID     *uint              `gorm:"default:1" json:"tier_id"`
	CreatedAt  time.Time          `json:"created_at"`
	UpdatedAt  time.Time          `json:"updated_at"`
	Sessions   []UserQuerySession `gorm:"foreignKey:UserID" json:"sessions"`
	Tier       *Tier              `gorm:"foreignKey:TierID" json:"tier"`
	MFA        *UserMFA           `gorm:"foreignKey:UserID" json:"mfa,omitempty"`
}

func (u User) GetId() uint {
	return u.ID
}

func (u User) GetIdType() string {
	return "uint"
}

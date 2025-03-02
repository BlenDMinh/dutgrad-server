package entity

import "time"

type UserAuthCredential struct {
	ID           uint
	UserID       uint
	User         User      `json:"id" gorm:"foreignKey:UserID"`
	AuthType     string    `json:"auth_type"`     // e.g., "local", "google", "facebook"
	PasswordHash *string   `json:"password_hash"` // For traditional login
	ExternalID   *string   `json:"external_id"`   // For Google/Facebook login
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

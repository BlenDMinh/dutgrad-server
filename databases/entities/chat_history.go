package entities

import (
	"time"

	"gorm.io/datatypes"
)

type ChatHistory struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	SessionID uint           `json:"session_id" gorm:"not null;index"`
	Message   datatypes.JSON `json:"message" gorm:"type:jsonb;not null"`
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
}

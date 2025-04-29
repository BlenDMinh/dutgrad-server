package entities

import "time"

type Space struct {
	ID              uint               `gorm:"primaryKey" json:"id"`
	Name            string             `gorm:"type:varchar(255);not null" json:"name"`
	Description     string             `gorm:"type:text" json:"description"`
	PrivacyStatus   bool               `json:"privacy_status"`
	DocumentLimit   int                `json:"document_limit" gorm:"default:10"`
	FileSizeLimitKb int                `json:"file_size_limit_kb" gorm:"default:5120"`
	ApiCallLimit    int                `json:"api_call_limit" gorm:"default:100"`
	CreatedAt       time.Time          `json:"created_at"`
	UpdatedAt       time.Time          `json:"updated_at"`
	Documents       []Document         `gorm:"foreignKey:SpaceID" json:"documents"`
	Sessions        []UserQuerySession `gorm:"foreignKey:SpaceID" json:"sessions"`
}

func (s Space) GetIdType() string {
	return "uint"
}

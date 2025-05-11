package entities

import "time"

type Document struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	SpaceID          uint      `gorm:"not null;index" json:"space_id"`
	Name             string    `gorm:"type:varchar(255)" json:"name"`
	Description      string    `gorm:"type:varchar(1024)" json:"description"`
	MimeType         string    `gorm:"column:mime_type;type:varchar(255)" json:"mime_type"`
	Size             int64     `gorm:"column:size;not null" json:"size"`
	ProcessingStatus int       `gorm:"default:0" json:"processing_status"`
	S3URL            string    `gorm:"not null" json:"s3_url"`
	PrivacyStatus    bool      `gorm:"default:true" json:"privacy_status"`
	CreatedAt        time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	Space            *Space    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"space"`
}

func (s Document) GetIdType() string {
	return "uint"
}

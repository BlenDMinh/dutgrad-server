package entities

import "time"

type UserQuerySession struct {
	ID            uint          `json:"id" gorm:"primaryKey"`
	UserID        uint          `json:"user_id" gorm:"not null;index"`
	SpaceID       uint          `json:"space_id" gorm:"not null;index"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
	User          *User         `json:"user" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`
	Space         *Space        `json:"space" gorm:"foreignKey:SpaceID;constraint:OnUpdate:CASCADE, OnDelete:CASCADE;"`
	UserQuery     []UserQuery   `json:"user_query" gorm:"foreignKey:QuerySessionID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ChatHistories []ChatHistory `json:"chat_histories" gorm:"foreignKey:SessionID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (u UserQuerySession) GetIdType() string {
	return "uint"
}

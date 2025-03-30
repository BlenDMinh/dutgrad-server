package entities

import "time"

type UserQuery struct {
	ID               uint             `json:"id" gorm:"primaryKey"`
	QuerySessionID   uint             `json:"query_session_id" gorm:"not null;index"`
	Query            string           `json:"query" gorm:"not null"`
	CreatedAt        time.Time        `json:"created_at"`
	UpdatedAt        time.Time        `json:"updated_at"`
	UserQuerySession UserQuerySession `json:"user_query_session" gorm:"foreignKey:QuerySessionID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (u UserQuery) GetIdType() string {
	return "uint"
}

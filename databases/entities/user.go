package entities

import "time"

type User struct {
	ID        uint               `json:"id" gorm:"primaryKey"`
	Username  string             `json:"username"`
	Email     *string            `json:"email" gorm:"unique"`
	Active    bool               `json:"active"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
	Sessions  []UserQuerySession `gorm:"foreignKey:UserID" json:"sessions"`
}

func (u User) GetIdType() string {
	return "uint"
}

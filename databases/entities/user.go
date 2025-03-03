package entities

import "time"

type User struct {
	ID          uint       `json:"id" gorm:"primaryKey"`
	FirstName   string     `json:"first_name"`
	LastName    string     `json:"last_name"`
	Email       *string    `json:"email" gorm:"unique"`
	Birthday    *time.Time `json:"birthday"`
	ActivatedAt *time.Time `json:"activated_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

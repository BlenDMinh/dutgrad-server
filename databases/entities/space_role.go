package entities

import "time"

type SpaceRole struct {
	ID         uint      `json:"id"`
	Name       string    `json:"name"`
	Permission uint      `json:"permission"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

package entities

import "time"

type SpaceRole struct {
	ID         uint      `json:"id"`
	Name       string    `json:"name"`
	Permission uint      `json:"permission"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

const (
	Owner  = 1
	Editor = 2
	Viewer = 3
)

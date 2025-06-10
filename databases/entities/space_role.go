package entities

import "time"

type SpaceRole struct {
	ID         uint      `json:"id"`
	Name       string    `json:"name"`
	Permission uint      `json:"permission"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (r SpaceRole) IsOwner() bool {
	return r.ID == SpaceRoleOwner
}

func (r SpaceRole) IsEditor() bool {
	return r.ID == SpaceRoleEditor
}

func (r SpaceRole) IsViewer() bool {
	return r.ID == SpaceRoleViewer
}

const (
	SpaceRoleOwner  = 1
	SpaceRoleEditor = 2
	SpaceRoleViewer = 3
)

package dtos

import "github.com/BlenDMinh/dutgrad-server/databases/entities"

type UserListResponse struct {
	Users []entities.User `json:"users"`
}

type UserTierUsageResponse struct {
	Tier  *entities.Tier `json:"tier"`
	Usage *TierUsage     `json:"usage"`
}

package dtos

import "github.com/BlenDMinh/dutgrad-server/databases/entities"

type TierUsageResponse struct {
	Tier  *entities.Tier `json:"tier"`
	Usage *TierUsage     `json:"usage"`
}

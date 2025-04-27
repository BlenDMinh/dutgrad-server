package dtos

import "github.com/BlenDMinh/dutgrad-server/databases/entities"

type TierUsageResponse struct {
	Tier  *entities.Tier `json:"tier"`
	Usage *TierUsage     `json:"usage"`
}

type TierUsage struct {
	SpaceCount        int64 `json:"space_count"`
	DocumentCount     int64 `json:"document_count"`
	QueryHistoryCount int64 `json:"query_history_count"`
	TodayQueryCount   int64 `json:"today_query_count"`
	TodayApiCallCount int64 `json:"today_api_call_count"`
}

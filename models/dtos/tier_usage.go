package dtos

type TierUsage struct {
	SpaceCount        int64 `json:"space_count"`
	DocumentCount     int64 `json:"document_count"`
	TotalChatMessages int64 `json:"total_chat_messages"`
	ChatUsageDaily    int64 `json:"chat_usage_daily"`
	ChatUsageMonthly  int64 `json:"chat_usage_monthly"`
}

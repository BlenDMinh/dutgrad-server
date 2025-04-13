package dtos

type AskRequest struct {
	QuerySessionID uint   `json:"query_session_id" binding:"required"`
	Query          string `json:"query" binding:"required"`
}

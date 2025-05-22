package dtos

type OAuthStateResponse struct {
	StateToken string `json:"state_token"`
	ExpiresAt  string `json:"expires_at"`
}

type OAuthErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type OAuthMFAResponse struct {
	TempToken string `json:"temp_token"`
	UserID    uint   `json:"user_id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
}

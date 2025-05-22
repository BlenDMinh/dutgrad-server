package dtos

// OAuthStateResponse represents the response when a state token is provided
type OAuthStateResponse struct {
	StateToken string `json:"state_token"`
	ExpiresAt  string `json:"expires_at"`
}

// OAuthErrorResponse represents an OAuth error response
type OAuthErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// OAuthMFAResponse represents an MFA response during OAuth
type OAuthMFAResponse struct {
	TempToken string `json:"temp_token"`
	UserID    uint   `json:"user_id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
}

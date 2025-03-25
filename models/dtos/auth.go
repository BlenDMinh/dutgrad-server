package dtos

type RegisterDTO struct {
	Username string
	Email    string
	Password string
}

type ExternalAuthDTO struct {
	Email      string
	Username   string
	ExternalID string
	AuthType   string
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type ExternalAuthRequest struct {
	TokenID    string `json:"token_id" binding:"required"`
	Email      string `json:"email" binding:"required,email"`
	Username   string `json:"username" binding:"required"`
	ExternalID string `json:"external_id" binding:"required"`
	AuthType   string `json:"auth_type" binding:"required"` // "google" or "facebook"
}

type AuthResponse struct {
	Token     string      `json:"token"`
	IsNewUser bool        `json:"is_new_user"`
	User      interface{} `json:"user"`
	Expires   interface{} `json:"expires"`
}

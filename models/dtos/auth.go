package dtos

type RegisterDTO struct {
	FirstName string
	LastName  string
	Email     string
	Password  string
}

type ExternalAuthDTO struct {
	Email      string
	FirstName  string
	LastName   string
	ExternalID string
	AuthType   string
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8"`
}

type ExternalAuthRequest struct {
	TokenID    string `json:"token_id" binding:"required"`
	Email      string `json:"email" binding:"required,email"`
	FirstName  string `json:"first_name" binding:"required"`
	LastName   string `json:"last_name" binding:"required"`
	ExternalID string `json:"external_id" binding:"required"`
	AuthType   string `json:"auth_type" binding:"required"` // "google" or "facebook"
}

type AuthResponse struct {
	Token   string      `json:"token"`
	User    interface{} `json:"user"`
	Expires interface{} `json:"expires"`
}

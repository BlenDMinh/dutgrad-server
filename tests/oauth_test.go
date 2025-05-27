package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupOAuthRouter() *gin.Engine {
	r := gin.Default()
	authGroup := r.Group("/auth")
	{
		oauthGroup := authGroup.Group("/oauth")
		oauthGroup.GET("/google", GoogleOAuthHandler)

		authGroup.POST("/external-auth", ExternalAuthHandler)
		authGroup.POST("/exchange-state", ExchangeStateHandler)
	}
	return r
}

// Mock OAuth states for security
var mockOAuthStates = map[string]bool{
	"valid-state": true,
}

// Mock external auth tokens
var mockExternalTokens = map[string]struct {
	UserID     uint
	Email      string
	AuthType   string
	ExternalID string
}{
	"valid-token": {
		UserID:     4,
		Email:      "google@example.com",
		AuthType:   "google",
		ExternalID: "123456789",
	},
}

func GoogleOAuthHandler(c *gin.Context) {
	state := "valid-state"
	redirectURI := "http://localhost:3000/auth/callback"
	authURL := "https://accounts.google.com/o/oauth2/v2/auth?client_id=mock-client-id&response_type=code&scope=email%20profile&state=" + state + "&redirect_uri=" + redirectURI

	c.Redirect(http.StatusTemporaryRedirect, authURL)
}

func ExternalAuthHandler(c *gin.Context) {
	var req struct {
		Token      string `json:"token"`
		AuthType   string `json:"auth_type"`
		ExternalID string `json:"external_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if req.AuthType != "google" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported auth type"})
		return
	}

	// Validate external token
	userData, exists := mockExternalTokens[req.Token]
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": map[string]interface{}{
			"id":    userData.UserID,
			"email": userData.Email,
		},
		"token": "jwt-token-" + userData.ExternalID,
	})
}

func ExchangeStateHandler(c *gin.Context) {
	var req struct {
		State string `json:"state"`
		Code  string `json:"code"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if !mockOAuthStates[req.State] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid state"})
		return
	}

	// Mock exchanging code for token
	if req.Code == "valid-code" {
		c.JSON(http.StatusOK, gin.H{
			"token": "valid-token",
		})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid code"})
	}
}

func TestGoogleOAuth(t *testing.T) {
	router := setupOAuthRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/auth/oauth/google", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusTemporaryRedirect, w.Code)
	location := w.Header().Get("Location")
	assert.Contains(t, location, "accounts.google.com")
	assert.Contains(t, location, "state=valid-state")
}

func TestExchangeState(t *testing.T) {
	router := setupOAuthRouter()

	tests := []struct {
		name         string
		state        string
		code         string
		expectedCode int
	}{
		{
			name:         "✅ Exchange thành công",
			state:        "valid-state",
			code:         "valid-code",
			expectedCode: http.StatusOK,
		},
		{
			name:         "❌ State không hợp lệ",
			state:        "invalid-state",
			code:         "valid-code",
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "❌ Code không hợp lệ",
			state:        "valid-state",
			code:         "invalid-code",
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqBody, _ := json.Marshal(map[string]string{
				"state": tt.state,
				"code":  tt.code,
			})

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/auth/exchange-state", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
		})
	}
}

func TestExternalAuth(t *testing.T) {
	router := setupOAuthRouter()

	tests := []struct {
		name         string
		token        string
		authType     string
		expectedCode int
	}{
		{
			name:         "✅ External auth thành công",
			token:        "valid-token",
			authType:     "google",
			expectedCode: http.StatusOK,
		},
		{
			name:         "❌ Token không hợp lệ",
			token:        "invalid-token",
			authType:     "google",
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "❌ Auth type không được hỗ trợ",
			token:        "valid-token",
			authType:     "facebook",
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqBody, _ := json.Marshal(map[string]string{
				"token":     tt.token,
				"auth_type": tt.authType,
			})

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/auth/external-auth", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)

			if tt.expectedCode == http.StatusOK {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "user")
				assert.Contains(t, response, "token")
			}
		})
	}
}

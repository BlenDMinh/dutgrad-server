package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// Mock router để test API
func setupLoginRouter() *gin.Engine {
	r := gin.Default()
	authGroup := r.Group("/auth")
	{
		authGroup.POST("/login", LoginHandler)
	}
	return r
}

type User struct {
	ID       uint   `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"-"`
}

// Mock database cho người dùng
var mockUsersDB = map[string]User{
	"user@example.com": {
		ID:       1,
		Email:    "user@example.com",
		Password: "password123",
	},
	"mfa@example.com": {
		ID:       3,
		Email:    "mfa@example.com",
		Password: "mfapassword",
	},
	"google@example.com": {
		ID:       4,
		Email:    "google@example.com",
		Password: "", // Người dùng Google không có mật khẩu
	},
}

var mockMFAUsers = map[uint]bool{
	3: true, // MFA
}

func checkLoginCredentials(email, password string) (*User, bool, error) {
	user, exists := mockUsersDB[email]
	if !exists {
		return nil, false, &LoginError{message: "user not found"}
	}

	//đăng nhập bên thứ ba (Google)
	if email == "google@example.com" {
		//không cần mật khẩu
		requiresMFA := mockMFAUsers[user.ID]
		return &user, requiresMFA, nil
	}

	// Kiểm tra mật khẩu cho đăng nhập thông thường
	if user.Password != password {
		return nil, false, &LoginError{message: "invalid password"}
	}

	requiresMFA := mockMFAUsers[user.ID]
	return &user, requiresMFA, nil
}

type LoginError struct {
	message string
}

func (e *LoginError) Error() string {
	return e.message
}

func LoginHandler(c *gin.Context) {
	var loginRequest struct {
		Email      string `json:"email"`
		Password   string `json:"password"`
		AuthType   string `json:"auth_type,omitempty"` // (google, facebook,...)
		ExternalID string `json:"external_id,omitempty"`
	}

	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	if loginRequest.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing email"})
		return
	}

	//đăng nhập bên thứ ba (Google)
	if loginRequest.AuthType == "google" && loginRequest.ExternalID != "" {
	} else if loginRequest.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing password"})
		return
	}

	user, requiresMFA, err := checkLoginCredentials(loginRequest.Email, loginRequest.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
		return
	}

	// không cần MFA, đăng nhập thành công
	if !requiresMFA {
		mockToken := "mock-jwt-token-for-" + user.Email
		expiresAt := time.Now().Add(24 * time.Hour)
		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"message": "Login successful",
			"data": gin.H{
				"token": mockToken,
				"user": gin.H{
					"id":       user.ID,
					"email":    user.Email,
					"username": user.Username,
				},
				"expires": expiresAt,
			},
		})
		return
	}
	// MFA, tạo token
	tempToken := "temp-token-for-" + user.Email
	expiresAt := time.Now().Add(10 * time.Minute)

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "MFA verification required",
		"data": gin.H{
			"requires_mfa": true,
			"temp_token":   tempToken,
			"expires_at":   expiresAt.Format(time.RFC3339),
			"user": gin.H{
				"id":       user.ID,
				"username": user.Username,
				"email":    user.Email,
			},
		},
	})
}

// ----- Test Cases ----- //
func TestLoginAPI(t *testing.T) {
	router := setupLoginRouter()
	tests := []struct {
		name         string
		requestBody  map[string]interface{}
		expectedCode int
		expectedMsg  string
		checkMFA     bool
		checkToken   bool
	}{{
		name: "✅ Đăng nhập thành công",
		requestBody: map[string]interface{}{
			"email":    "user@example.com",
			"password": "password123",
		},
		expectedCode: http.StatusOK,
		expectedMsg:  "Login successful",
		checkMFA:     false,
		checkToken:   true}, {
		name: "✅ Đăng nhập Google thành công",
		requestBody: map[string]interface{}{
			"email":       "google@example.com",
			"auth_type":   "google",
			"external_id": "123456789",
		},
		expectedCode: http.StatusOK,
		expectedMsg:  "Login successful",
		checkMFA:     false,
		checkToken:   true,
	}, {
		name: "✅ Đăng nhập yêu cầu MFA",
		requestBody: map[string]interface{}{
			"email":    "mfa@example.com",
			"password": "mfapassword",
		},
		expectedCode: http.StatusOK,
		expectedMsg:  "MFA verification required",
		checkMFA:     true,
		checkToken:   false,
	},
		{
			name: "❌ Email không tồn tại",
			requestBody: map[string]interface{}{
				"email":    "nonexistent@example.com",
				"password": "password123",
			},
			expectedCode: http.StatusUnauthorized,
			expectedMsg:  "Authentication failed",
			checkMFA:     false,
			checkToken:   false,
		},
		{
			name: "❌ Sai mật khẩu",
			requestBody: map[string]interface{}{
				"email":    "user@example.com",
				"password": "wrongpassword",
			},
			expectedCode: http.StatusUnauthorized,
			expectedMsg:  "Authentication failed",
			checkMFA:     false,
			checkToken:   false,
		},
		{
			name: "❌ Thiếu email",
			requestBody: map[string]interface{}{
				"password": "password123",
			},
			expectedCode: http.StatusBadRequest,
			expectedMsg:  "Missing email",
			checkMFA:     false,
			checkToken:   false,
		}, {
			name: "❌ Thiếu mật khẩu",
			requestBody: map[string]interface{}{
				"email": "user@example.com",
			},
			expectedCode: http.StatusBadRequest,
			expectedMsg:  "Missing password",
			checkMFA:     false,
			checkToken:   false,
		}, {
			name: "❌ Thiếu auth_type khi đăng nhập Google",
			requestBody: map[string]interface{}{
				"email":       "google@example.com",
				"external_id": "987654321",
			},
			expectedCode: http.StatusBadRequest,
			expectedMsg:  "Missing password",
			checkMFA:     false,
			checkToken:   false,
		},
		{
			name:         "❌ Body rỗng",
			requestBody:  map[string]interface{}{},
			expectedCode: http.StatusBadRequest,
			expectedMsg:  "Missing email",
			checkMFA:     false,
			checkToken:   false,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			body, _ := json.Marshal(tc.requestBody)
			req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedCode, w.Code)

			var response map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &response)

			if tc.expectedCode == http.StatusOK {
				assert.Equal(t, float64(tc.expectedCode), response["status"])
				assert.Equal(t, tc.expectedMsg, response["message"])
				data, _ := response["data"].(map[string]interface{})

				if tc.checkMFA {
					requiresMFA, exists := data["requires_mfa"].(bool)
					assert.True(t, exists)
					assert.True(t, requiresMFA)
					assert.NotNil(t, data["temp_token"])
					assert.NotNil(t, data["expires_at"])
				}

				if tc.checkToken {
					assert.NotNil(t, data["token"])
					assert.NotNil(t, data["expires"])
					assert.NotNil(t, data["user"])
				}
			} else {
				assert.Equal(t, tc.expectedMsg, response["error"])
			}
		})
	}
}

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

// Mock router để test API
func setupRouter() *gin.Engine {
	r := gin.Default()
	r.POST("/register", RegisterHandler)
	return r
}

// Mock database để kiểm tra email tồn tại
var mockUsers = map[string]bool{
	"existing@example.com": true,
}

func isEmailExists(email string) bool {
	return mockUsers[email]
}

func RegisterHandler(c *gin.Context) {
	var user struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
		Password  string `json:"password"`
		Birthday  string `json:"birthday"`
	}

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if user.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing email"})
		return
	}
	if !isValidEmail(user.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email"})
		return
	}
	if isEmailExists(user.Email) {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
		return
	}

	if user.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing password"})
		return
	}
	if len(user.Password) < 8 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password too short"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

func isValidEmail(email string) bool {
	return len(email) > 3 && len(email) < 50 && contains(email, "@")
}

func contains(str, substr string) bool {
	return bytes.Contains([]byte(str), []byte(substr))
}

// ----- Test Cases ----- //
func TestRegisterAPI(t *testing.T) {
	router := setupRouter()

	tests := []struct {
		name         string
		requestBody  map[string]string
		expectedCode int
		expectedMsg  string
	}{
		{
			name: "✅ Đăng ký thành công",
			requestBody: map[string]string{
				"first_name": "John",
				"last_name":  "Doe",
				"email":      "newuser@example.com",
				"password":   "12345678",
			},
			expectedCode: http.StatusCreated,
			expectedMsg:  "User registered successfully",
		},
		{
			name: "❌ Email đã tồn tại",
			requestBody: map[string]string{
				"first_name": "Jane",
				"last_name":  "Doe",
				"email":      "existing@example.com",
				"password":   "12345678",
			},
			expectedCode: http.StatusConflict,
			expectedMsg:  "Email already exists",
		},
		{
			name: "❌ Thiếu email",
			requestBody: map[string]string{
				"first_name": "Mike",
				"last_name":  "Smith",
				"password":   "12345678",
			},
			expectedCode: http.StatusBadRequest,
			expectedMsg:  "Missing email",
		},
		{
			name: "❌ Thiếu trường password",
			requestBody: map[string]string{
				"first_name": "John",
				"last_name":  "Doe",
				"email":      "test@example.com",
			},
			expectedCode: http.StatusBadRequest,
			expectedMsg:  "Missing password",
		},
		{
			name: "❌ Thiếu trường email và password",
			requestBody: map[string]string{
				"first_name": "John",
				"last_name":  "Doe",
			},
			expectedCode: http.StatusBadRequest,
			expectedMsg:  "Missing email",
		},
		{
			name: "❌ Email không hợp lệ",
			requestBody: map[string]string{
				"first_name": "Lucas",
				"last_name":  "Brown",
				"email":      "invalid-email",
				"password":   "12345678",
			},
			expectedCode: http.StatusBadRequest,
			expectedMsg:  "Invalid email",
		},
		{
			name: "❌ Mật khẩu ngắn",
			requestBody: map[string]string{
				"first_name": "Emma",
				"last_name":  "Taylor",
				"email":      "emma@example.com",
				"password":   "123",
			},
			expectedCode: http.StatusBadRequest,
			expectedMsg:  "Password too short",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			body, _ := json.Marshal(tc.requestBody)
			req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedCode, w.Code)

			var response map[string]string
			json.Unmarshal(w.Body.Bytes(), &response)
			if tc.expectedCode == http.StatusCreated {
				assert.Equal(t, tc.expectedMsg, response["message"])
			} else {
				assert.Equal(t, tc.expectedMsg, response["error"])
			}
		})
	}
}

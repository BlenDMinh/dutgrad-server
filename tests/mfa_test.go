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

func setupMFARouter() *gin.Engine {
	r := gin.Default()
	authGroup := r.Group("/auth/mfa")
	{
		authGroup.GET("/status", MFAStatusHandler)
		authGroup.POST("/setup", SetupMFAHandler)
		authGroup.POST("/verify", VerifyMFAHandler)
		authGroup.POST("/disable", DisableMFAHandler)
	}
	return r
}

// Mock user MFA data
var mockMFAData = map[uint]struct {
	Secret    string
	IsEnabled bool
}{
	1: {Secret: "TESTSECRET123", IsEnabled: false},
	3: {Secret: "MFASECRET456", IsEnabled: true},
}

func MFAStatusHandler(c *gin.Context) {
	userID := uint(1) // Mock authenticated user

	mfaData, exists := mockMFAData[userID]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "User MFA data not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"is_enabled": mfaData.IsEnabled,
	})
}

func SetupMFAHandler(c *gin.Context) {
	userID := uint(1) // Mock authenticated user

	// Generate new MFA secret
	newSecret := "NEWSECRET789"

	mockMFAData[userID] = struct {
		Secret    string
		IsEnabled bool
	}{
		Secret:    newSecret,
		IsEnabled: false,
	}

	c.JSON(http.StatusOK, gin.H{
		"secret":      newSecret,
		"qr_code_url": "otpauth://totp/DUTGrad:user@example.com?secret=" + newSecret + "&issuer=DUTGrad",
	})
}

func VerifyMFAHandler(c *gin.Context) {
	var req struct {
		Code string `json:"code"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	userID := uint(1) // Mock authenticated user
	mfaData, exists := mockMFAData[userID]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "User MFA data not found"})
		return
	}

	// Simulate TOTP verification
	if req.Code == "123456" { // Mock valid code
		mockMFAData[userID] = struct {
			Secret    string
			IsEnabled bool
		}{
			Secret:    mfaData.Secret,
			IsEnabled: true,
		}
		c.JSON(http.StatusOK, gin.H{"message": "MFA verified and enabled"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid MFA code"})
	}
}

func DisableMFAHandler(c *gin.Context) {
	userID := uint(1) // Mock authenticated user

	mfaData, exists := mockMFAData[userID]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "User MFA data not found"})
		return
	}

	if !mfaData.IsEnabled {
		c.JSON(http.StatusBadRequest, gin.H{"error": "MFA is not enabled"})
		return
	}

	mockMFAData[userID] = struct {
		Secret    string
		IsEnabled bool
	}{
		Secret:    mfaData.Secret,
		IsEnabled: false,
	}

	c.JSON(http.StatusOK, gin.H{"message": "MFA disabled successfully"})
}

func TestMFAStatus(t *testing.T) {
	router := setupMFARouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/auth/mfa/status", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response, "is_enabled")
}

func TestSetupMFA(t *testing.T) {
	router := setupMFARouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/auth/mfa/setup", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response, "secret")
	assert.Contains(t, response, "qr_code_url")
}

func TestVerifyMFA(t *testing.T) {
	router := setupMFARouter()

	tests := []struct {
		name         string
		code         string
		expectedCode int
	}{
		{
			name:         "✅ Verify MFA thành công",
			code:         "123456",
			expectedCode: http.StatusOK,
		},
		{
			name:         "❌ Code không hợp lệ",
			code:         "111111",
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqBody, _ := json.Marshal(map[string]string{
				"code": tt.code,
			})

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/auth/mfa/verify", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
		})
	}
}

func TestDisableMFA(t *testing.T) {
	router := setupMFARouter()

	// Enable MFA first
	mockMFAData[1] = struct {
		Secret    string
		IsEnabled bool
	}{
		Secret:    "TESTSECRET123",
		IsEnabled: true,
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/auth/mfa/disable", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "MFA disabled successfully", response["message"])
}

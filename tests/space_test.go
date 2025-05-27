package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type Space struct {
	ID              uint      `json:"id"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	PrivacyStatus   bool      `json:"privacy_status"`
	SystemPrompt    string    `json:"system_prompt"`
	DocumentLimit   int       `json:"document_limit"`
	FileSizeLimitKb int       `json:"file_size_limit_kb"`
	ApiCallLimit    int       `json:"api_call_limit"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type SpaceUser struct {
	UserID      uint  `json:"user_id"`
	SpaceID     uint  `json:"space_id"`
	SpaceRoleID *uint `json:"space_role_id"`
}

type SpaceRole struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Permission  uint      `json:"permission"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

var mockSpaces = []Space{
	{
		ID:              1,
		Name:            "Public Space",
		Description:     "This is a public space",
		PrivacyStatus:   false,
		SystemPrompt:    "Default system prompt",
		DocumentLimit:   10,
		FileSizeLimitKb: 5120,
		ApiCallLimit:    100,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	},
	{
		ID:              2,
		Name:            "Private Space",
		Description:     "This is a private space",
		PrivacyStatus:   true,
		SystemPrompt:    "Default system prompt",
		DocumentLimit:   10,
		FileSizeLimitKb: 5120,
		ApiCallLimit:    100,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	},
}

var mockSpaceUsers = []SpaceUser{
	{
		UserID:      1,
		SpaceID:     1,
		SpaceRoleID: pointerToUint(1), // Owner
	},
	{
		UserID:      2,
		SpaceID:     1,
		SpaceRoleID: pointerToUint(2), // Editor
	},
	{
		UserID:      1,
		SpaceID:     2,
		SpaceRoleID: pointerToUint(1), // Owner
	},
}

// Mock data for space roles
var mockSpaceRoles = []SpaceRole{
	{
		ID:          1,
		Name:        "Owner",
		Description: "Owner of the space",
		Permission:  Owner,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	},
	{
		ID:          2,
		Name:        "Editor",
		Description: "Editor of the space",
		Permission:  Editor,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	},
	{
		ID:          3,
		Name:        "Viewer",
		Description: "Viewer of the space",
		Permission:  Viewer,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	},
}

const (
	Owner  = 1
	Editor = 2
	Viewer = 3
)

func pointerToUint(i uint) *uint {
	return &i
}

func setupSpaceRouter() *gin.Engine {
	r := gin.Default()

	r.Use(mockAuthMiddleware())

	r.GET("/spaces", GetPublicSpacesHandler)
	r.GET("/spaces/:id", GetSpaceHandler)
	r.GET("/spaces/roles", GetSpaceRolesHandler)
	r.GET("/spaces/popular", GetPopularSpacesHandler)
	r.GET("/spaces/me", GetMySpacesHandler)
	r.HEAD("/spaces/count/me", CountMySpacesHandler)
	r.GET("/spaces/user/:user_id", GetUserSpacesHandler)
	r.GET("/spaces/public", GetPublicSpacesFromPublicEndpointHandler)
	r.PATCH("/spaces/:id/members/:memberId/role", UpdateUserRoleHandler)
	r.DELETE("/spaces/:id/members/:memberId", RemoveMemberHandler)
	r.POST("/spaces/:id/chat", ChatHandler)

	r.POST("/spaces", CreateSpaceHandler)
	r.PUT("/spaces/:id", UpdateSpaceHandler)
	r.DELETE("/spaces/:id", DeleteSpaceHandler)
	r.GET("/spaces/:id/members", GetSpaceMembersHandler)
	r.GET("/spaces/:id/invitations", GetSpaceInvitationsHandler)
	r.PUT("/spaces/:id/invitation-link", GetInvitationLinkHandler)
	r.POST("/spaces/:id/invitations", InviteUserToSpaceHandler)
	r.GET("/spaces/:id/user-role", GetUserRoleHandler)
	r.POST("/spaces/join", JoinSpaceHandler)
	r.POST("/spaces/:id/join-public", JoinPublicSpaceHandler)

	r.GET("/spaces/:id/api-keys", ListApiKeysHandler)
	r.POST("/spaces/:id/api-keys", CreateApiKeyHandler)
	r.GET("/spaces/:id/api-keys/:keyId", GetApiKeyHandler)
	r.DELETE("/spaces/:id/api-keys/:keyId", DeleteApiKeyHandler)

	return r
}

func mockAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	}
}

// GET /spaces - Get public spaces
func GetPublicSpacesHandler(c *gin.Context) {
	var publicSpaces []Space

	for _, space := range mockSpaces {
		if !space.PrivacyStatus {
			publicSpaces = append(publicSpaces, space)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Success",
		"data": gin.H{
			"public_spaces": publicSpaces,
			"pagination": gin.H{
				"current_page": 1,
				"page_size":    10,
				"total_pages":  1,
				"total_items":  len(publicSpaces),
				"has_next":     false,
				"has_prev":     false,
			},
		},
	})
}

// GET /spaces/public - Get public spaces from public endpoint
func GetPublicSpacesFromPublicEndpointHandler(c *gin.Context) {
	var publicSpaces []Space

	for _, space := range mockSpaces {
		if !space.PrivacyStatus {
			publicSpaces = append(publicSpaces, space)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Success",
		"data": gin.H{
			"spaces": publicSpaces,
			"pagination": gin.H{
				"current_page": 1,
				"page_size":    10,
				"total_pages":  1,
				"total_items":  len(publicSpaces),
				"has_next":     false,
				"has_prev":     false,
			},
		},
	})
}

// GET /spaces/:id - Get a specific space
func GetSpaceHandler(c *gin.Context) {
	spaceID := c.Param("id")

	var foundSpace *Space
	for _, space := range mockSpaces {
		if fmt.Sprintf("%d", space.ID) == spaceID {
			foundSpace = &space
			break
		}
	}

	if foundSpace == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "Space not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Success",
		"data":    foundSpace,
	})
}

// POST /spaces - Create a new space
func CreateSpaceHandler(c *gin.Context) {
	var newSpace Space

	if err := c.ShouldBindJSON(&newSpace); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Invalid request",
		})
		return
	}

	// Set default values and ID
	newSpace.ID = uint(len(mockSpaces) + 1)
	newSpace.CreatedAt = time.Now()
	newSpace.UpdatedAt = time.Now()

	if newSpace.SystemPrompt == "" {
		newSpace.SystemPrompt = "You are an AI assistant for answering questions about documents in this space. Provide helpful, accurate, and concise information based on the content available."
	}

	if newSpace.DocumentLimit == 0 {
		newSpace.DocumentLimit = 10
	}

	if newSpace.FileSizeLimitKb == 0 {
		newSpace.FileSizeLimitKb = 5120
	}

	if newSpace.ApiCallLimit == 0 {
		newSpace.ApiCallLimit = 100
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  http.StatusCreated,
		"message": "Created",
		"data":    newSpace,
	})
}

// PUT /spaces/:id - Update a space
func UpdateSpaceHandler(c *gin.Context) {
	spaceID := c.Param("id")
	var updatedSpace Space

	if err := c.ShouldBindJSON(&updatedSpace); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	found := false
	for i, space := range mockSpaces {
		if fmt.Sprintf("%d", space.ID) == spaceID {
			updatedSpace.ID = space.ID
			updatedSpace.CreatedAt = space.CreatedAt
			updatedSpace.UpdatedAt = time.Now()
			mockSpaces[i] = updatedSpace
			found = true
			break
		}
	}

	if !found {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "Space not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Updated",
		"data":    updatedSpace,
	})
}

// DELETE /spaces/:id - Delete a space
func DeleteSpaceHandler(c *gin.Context) {
	spaceID := c.Param("id")

	found := false
	for i, space := range mockSpaces {
		if fmt.Sprintf("%d", space.ID) == spaceID {
			mockSpaces = append(mockSpaces[:i], mockSpaces[i+1:]...)
			found = true
			break
		}
	}

	if !found {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "Space not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Deleted",
	})
}

// GET /spaces/:id/members - Get space members
func GetSpaceMembersHandler(c *gin.Context) {
	spaceID := c.Param("id")
	spaceIDInt := 0
	fmt.Sscanf(spaceID, "%d", &spaceIDInt)

	var members []SpaceUser
	for _, spaceUser := range mockSpaceUsers {
		if spaceUser.SpaceID == uint(spaceIDInt) {
			members = append(members, spaceUser)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Success",
		"data": gin.H{
			"members": members,
		},
	})
}

// GET /spaces/:id/invitations - Get space invitations
func GetSpaceInvitationsHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Success",
		"data": gin.H{
			"invitations": []interface{}{},
		},
	})
}

// PUT /spaces/:id/invitation-link - Get invitation link
func GetInvitationLinkHandler(c *gin.Context) {
	spaceID := c.Param("id")

	invitationLink := fmt.Sprintf("https://example.com/invitation?token=mock_token_for_space_%s", spaceID)

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Success",
		"data": gin.H{
			"invitation_link": invitationLink,
		},
	})
}

// POST /spaces/:id/invitations - Invite a user to space
func InviteUserToSpaceHandler(c *gin.Context) {
	var req struct {
		InvitedUserID    *uint  `json:"invited_user_id"`
		InvitedUserEmail string `json:"invited_user_email"`
		SpaceRoleID      uint   `json:"space_role_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Invalid request",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Invitation sent successfully",
	})
}

// GET /spaces/:id/user-role - Get user role in space
func GetUserRoleHandler(c *gin.Context) {
	spaceID := c.Param("id")
	userID, _ := c.Get("user_id")
	spaceIDInt := 0
	fmt.Sscanf(spaceID, "%d", &spaceIDInt)
	var role *SpaceRole
	for _, spaceUser := range mockSpaceUsers {
		if spaceUser.SpaceID == uint(spaceIDInt) && spaceUser.UserID == userID.(uint) {
			for _, spaceRole := range mockSpaceRoles {
				if spaceRole.ID == *spaceUser.SpaceRoleID {
					role = &spaceRole
					break
				}
			}
			break
		}
	}

	if role == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "Role not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Success",
		"data": gin.H{
			"role": role,
		},
	})
}

// GET /spaces/roles - Get all space roles
func GetSpaceRolesHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Success",
		"data":    mockSpaceRoles,
	})
}

// POST /spaces/join - Join a space with invitation token
func JoinSpaceHandler(c *gin.Context) {
	token := c.Query("token")

	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Token is required",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Successfully joined the space",
		"data": gin.H{
			"space_id": 1,
		},
	})
}

// POST /spaces/:id/join-public - Join a public space
func JoinPublicSpaceHandler(c *gin.Context) {
	spaceID := c.Param("id")
	spaceIDInt := 0
	fmt.Sscanf(spaceID, "%d", &spaceIDInt)

	found := false
	isPublic := false
	for _, space := range mockSpaces {
		if space.ID == uint(spaceIDInt) {
			found = true
			if !space.PrivacyStatus {
				isPublic = true
			}
			break
		}
	}

	if !found {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "Space not found",
		})
		return
	}

	if !isPublic {
		c.JSON(http.StatusForbidden, gin.H{
			"status":  http.StatusForbidden,
			"message": "Space is not public",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Successfully joined the space",
	})
}

// GET /spaces/popular - Get popular spaces
func GetPopularSpacesHandler(c *gin.Context) {
	var publicSpaces []Space
	for _, space := range mockSpaces {
		if !space.PrivacyStatus {
			publicSpaces = append(publicSpaces, space)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Success",
		"data":    publicSpaces,
	})
}

// GET /spaces/me - Get spaces owned by current user
func GetMySpacesHandler(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var userSpaces []Space
	for _, spaceUser := range mockSpaceUsers {
		if spaceUser.UserID == userID.(uint) {
			for _, space := range mockSpaces {
				if space.ID == spaceUser.SpaceID {
					userSpaces = append(userSpaces, space)
					break
				}
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Success",
		"data": gin.H{
			"spaces": userSpaces,
			"pagination": gin.H{
				"current_page": 1,
				"page_size":    10,
				"total_pages":  1,
				"total_items":  len(userSpaces),
				"has_next":     false,
				"has_prev":     false,
			},
		},
	})
}

// HEAD /spaces/count/me - Count spaces owned by current user
func CountMySpacesHandler(c *gin.Context) {
	userID, _ := c.Get("user_id")

	count := 0
	for _, spaceUser := range mockSpaceUsers {
		if spaceUser.UserID == userID.(uint) {
			count++
		}
	}

	c.Header("X-Total-Count", fmt.Sprintf("%d", count))
	c.Status(http.StatusOK)
}

// GET /spaces/user/:user_id - Get spaces owned by a specific user
func GetUserSpacesHandler(c *gin.Context) {
	userIDParam := c.Param("user_id")
	userIDInt := uint(0)
	fmt.Sscanf(userIDParam, "%d", &userIDInt)

	var userSpaces []Space
	for _, spaceUser := range mockSpaceUsers {
		if spaceUser.UserID == userIDInt {
			for _, space := range mockSpaces {
				if space.ID == spaceUser.SpaceID && !space.PrivacyStatus {
					userSpaces = append(userSpaces, space)
					break
				}
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Success",
		"data": gin.H{
			"spaces": userSpaces,
			"pagination": gin.H{
				"current_page": 1,
				"page_size":    10,
				"total_pages":  1,
				"total_items":  len(userSpaces),
				"has_next":     false,
				"has_prev":     false,
			},
		},
	})
}

// PATCH /spaces/:id/members/:memberId/role - Update a user's role in a space
func UpdateUserRoleHandler(c *gin.Context) {
	spaceID := c.Param("id")
	memberID := c.Param("memberId")
	spaceIDInt := uint(0)
	memberIDInt := uint(0)
	fmt.Sscanf(spaceID, "%d", &spaceIDInt)
	fmt.Sscanf(memberID, "%d", &memberIDInt)

	var req struct {
		SpaceRoleID uint `json:"space_role_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Invalid request",
		})
		return
	}

	// Find the space user and update the role
	found := false
	for i, spaceUser := range mockSpaceUsers {
		if spaceUser.SpaceID == spaceIDInt && spaceUser.UserID == memberIDInt {
			newRoleID := req.SpaceRoleID
			mockSpaceUsers[i].SpaceRoleID = &newRoleID
			found = true
			break
		}
	}

	if !found {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "Member not found in space",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Role updated successfully",
	})
}

// DELETE /spaces/:id/members/:memberId - Remove a member from a space
func RemoveMemberHandler(c *gin.Context) {
	spaceID := c.Param("id")
	memberID := c.Param("memberId")
	spaceIDInt := uint(0)
	memberIDInt := uint(0)
	fmt.Sscanf(spaceID, "%d", &spaceIDInt)
	fmt.Sscanf(memberID, "%d", &memberIDInt)

	// Find the space user and remove them
	found := false
	for i, spaceUser := range mockSpaceUsers {
		if spaceUser.SpaceID == spaceIDInt && spaceUser.UserID == memberIDInt {
			mockSpaceUsers = append(mockSpaceUsers[:i], mockSpaceUsers[i+1:]...)
			found = true
			break
		}
	}

	if !found {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "Member not found in space",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Member removed successfully",
	})
}

// POST /spaces/:id/chat - Chat with AI in a space
func ChatHandler(c *gin.Context) {
	var req struct {
		Message string `json:"message"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Invalid request",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Success",
		"data": gin.H{
			"response": "This is a mock AI response to: " + req.Message,
		},
	})
}

// GET spaces/:id/api-keys - List API keys for a space
func ListApiKeysHandler(c *gin.Context) {
	spaceID := c.Param("id")

	apiKeys := []map[string]interface{}{
		{
			"id":          1,
			"name":        "Default API Key",
			"key":         "sk-****************************************",
			"space_id":    spaceID,
			"created_at":  time.Now(),
			"updated_at":  time.Now(),
			"expires_at":  time.Now().AddDate(0, 1, 0), // 1 month from now
			"is_disabled": false,
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Success",
		"data":    apiKeys,
	})
}

// POST spaces/:id/api-keys - Create new API key for a space
func CreateApiKeyHandler(c *gin.Context) {
	spaceID := c.Param("id")

	var req struct {
		Name      string     `json:"name"`
		ExpiresAt *time.Time `json:"expires_at"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Invalid request",
		})
		return
	}

	// Generate a mock API key
	newApiKey := map[string]interface{}{
		"id":          1,
		"name":        req.Name,
		"key":         "sk-" + generateRandomString(40),
		"space_id":    spaceID,
		"created_at":  time.Now(),
		"updated_at":  time.Now(),
		"expires_at":  req.ExpiresAt,
		"is_disabled": false,
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  http.StatusCreated,
		"message": "API key created successfully",
		"data":    newApiKey,
	})
}

// GET spaces/:id/api-keys/:keyId - Get API key details
func GetApiKeyHandler(c *gin.Context) {
	spaceID := c.Param("id")
	keyID := c.Param("keyId")

	apiKey := map[string]interface{}{
		"id":          keyID,
		"name":        "Default API Key",
		"key":         "sk-****************************************",
		"space_id":    spaceID,
		"created_at":  time.Now(),
		"updated_at":  time.Now(),
		"expires_at":  time.Now().AddDate(0, 1, 0), // 1 month from now
		"is_disabled": false,
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Success",
		"data":    apiKey,
	})
}

// DELETE spaces/:id/api-keys/:keyId - Delete an API key
func DeleteApiKeyHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "API key deleted successfully",
	})
}

func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}

func TestGetPublicSpaces(t *testing.T) {
	router := setupSpaceRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/spaces", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
		Data    struct {
			PublicSpaces []Space `json:"public_spaces"`
		} `json:"data"`
	}

	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, "Success", response.Message)
	assert.NotEmpty(t, response.Data.PublicSpaces, "Should return public spaces")

	for _, space := range response.Data.PublicSpaces {
		assert.False(t, space.PrivacyStatus, "Only public spaces should be returned")
	}
}

func TestGetSpaceDetails(t *testing.T) {
	router := setupSpaceRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/spaces/1", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
		Data    Space  `json:"data"`
	}

	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, "Success", response.Message)
	assert.Equal(t, uint(1), response.Data.ID)
	assert.Equal(t, "Public Space", response.Data.Name)
}

func TestCreateSpace(t *testing.T) {
	router := setupSpaceRouter()

	newSpace := Space{
		Name:          "Test Space",
		Description:   "Test Description",
		PrivacyStatus: true,
	}

	payload, _ := json.Marshal(newSpace)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/spaces", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
		Data    Space  `json:"data"`
	}

	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, "Created", response.Message)
	assert.Equal(t, "Test Space", response.Data.Name)
	assert.Equal(t, "Test Description", response.Data.Description)
	assert.Equal(t, true, response.Data.PrivacyStatus)

	assert.NotEmpty(t, response.Data.SystemPrompt)
	assert.Equal(t, 10, response.Data.DocumentLimit)
	assert.Equal(t, 5120, response.Data.FileSizeLimitKb)
	assert.Equal(t, 100, response.Data.ApiCallLimit)
}

func TestGetSpaceMembers(t *testing.T) {
	router := setupSpaceRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/spaces/1/members", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
		Data    struct {
			Members []SpaceUser `json:"members"`
		} `json:"data"`
	}

	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, "Success", response.Message)

	// Space 1 should have 2 members
	assert.Equal(t, 2, len(response.Data.Members))
}

func TestGetSpaceRoles(t *testing.T) {
	router := setupSpaceRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/spaces/roles", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		Status  int         `json:"status"`
		Message string      `json:"message"`
		Data    []SpaceRole `json:"data"`
	}

	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, "Success", response.Message)
	// Should have all roles
	assert.Equal(t, 3, len(response.Data))

	// Check if Owner role exists
	foundOwner := false
	for _, role := range response.Data {
		if role.Name == "Owner" {
			foundOwner = true
			break
		}
	}
	assert.True(t, foundOwner, "Owner role should exist")
}

func TestJoinSpaceWithToken(t *testing.T) {
	router := setupSpaceRouter()

	// Test with valid token
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/spaces/join?token=valid_token", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
		Data    struct {
			SpaceID uint `json:"space_id"`
		} `json:"data"`
	}

	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, "Successfully joined the space", response.Message)
	assert.Equal(t, uint(1), response.Data.SpaceID)

	// Test without token
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/spaces/join", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestJoinPublicSpace(t *testing.T) {
	router := setupSpaceRouter()

	// Test joining a public space
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/spaces/1/join-public", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
	}

	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, "Successfully joined the space", response.Message)

	// Test joining a private space
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/spaces/2/join-public", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestGetInvitationLink(t *testing.T) {
	router := setupSpaceRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/spaces/1/invitation-link", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
		Data    struct {
			InvitationLink string `json:"invitation_link"`
		} `json:"data"`
	}

	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, "Success", response.Message)
	assert.Contains(t, response.Data.InvitationLink, "mock_token_for_space_1")
}

func TestInviteUserToSpace(t *testing.T) {
	router := setupSpaceRouter()

	invitation := struct {
		InvitedUserID    *uint  `json:"invited_user_id"`
		InvitedUserEmail string `json:"invited_user_email"`
		SpaceRoleID      uint   `json:"space_role_id"`
	}{
		InvitedUserEmail: "test@example.com",
		SpaceRoleID:      3,
	}

	payload, _ := json.Marshal(invitation)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/spaces/1/invitations", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
	}

	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, "Invitation sent successfully", response.Message)
}

func TestGetPopularSpaces(t *testing.T) {
	router := setupSpaceRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/spaces/popular", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		Status  int     `json:"status"`
		Message string  `json:"message"`
		Data    []Space `json:"data"`
	}

	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, "Success", response.Message)
	assert.NotEmpty(t, response.Data, "Should return public spaces")

	for _, space := range response.Data {
		assert.False(t, space.PrivacyStatus, "Only public spaces should be returned")
	}
}

func TestGetMySpaces(t *testing.T) {
	router := setupSpaceRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/spaces/me", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
		Data    struct {
			Spaces []Space `json:"spaces"`
		} `json:"data"`
	}

	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, "Success", response.Message)

	// User ID 1 should have 2 spaces (ID 1 and 2)
	assert.Equal(t, 2, len(response.Data.Spaces))

	// Verify that the spaces returned are the ones user 1 is a member of
	spaceIDs := make([]uint, 0)
	for _, space := range response.Data.Spaces {
		spaceIDs = append(spaceIDs, space.ID)
	}
	assert.Contains(t, spaceIDs, uint(1))
	assert.Contains(t, spaceIDs, uint(2))
}

func TestCountMySpaces(t *testing.T) {
	router := setupSpaceRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("HEAD", "/spaces/count/me", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// User ID 1 should have 2 spaces
	assert.Equal(t, "2", w.Header().Get("X-Total-Count"))
}

func TestUpdateUserRole(t *testing.T) {
	router := setupSpaceRouter()

	// Change user 2's role in space 1 from Editor to Viewer
	roleUpdate := struct {
		SpaceRoleID uint `json:"space_role_id"`
	}{
		SpaceRoleID: 3,
	}

	payload, _ := json.Marshal(roleUpdate)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/spaces/1/members/2/role", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
	}

	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, "Role updated successfully", response.Message)

	// Verify the role was updated
	var viewerRoleID uint = 3
	assert.Equal(t, &viewerRoleID, mockSpaceUsers[1].SpaceRoleID)
}

func TestRemoveMember(t *testing.T) {
	router := setupSpaceRouter()

	initialMembersCount := len(mockSpaceUsers)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/spaces/1/members/2", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
	}

	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, "Member removed successfully", response.Message)

	assert.Equal(t, initialMembersCount-1, len(mockSpaceUsers))

	// Verify user 2 is no longer in space 1
	userFound := false
	for _, spaceUser := range mockSpaceUsers {
		if spaceUser.SpaceID == 1 && spaceUser.UserID == 2 {
			userFound = true
			break
		}
	}
	assert.False(t, userFound, "User should be removed from the space")
}

func TestChatWithAI(t *testing.T) {
	router := setupSpaceRouter()

	chatMessage := struct {
		Message string `json:"message"`
	}{
		Message: "Hello AI!",
	}

	payload, _ := json.Marshal(chatMessage)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/spaces/1/chat", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
		Data    struct {
			Response string `json:"response"`
		} `json:"data"`
	}

	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, "Success", response.Message)
	assert.Contains(t, response.Data.Response, "Hello AI!")
}

func TestGetPublicSpacesFromPublicEndpoint(t *testing.T) {
	router := setupSpaceRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/spaces/public", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
		Data    struct {
			Spaces     []Space `json:"spaces"`
			Pagination struct {
				CurrentPage int  `json:"current_page"`
				PageSize    int  `json:"page_size"`
				TotalPages  int  `json:"total_pages"`
				TotalItems  int  `json:"total_items"`
				HasNext     bool `json:"has_next"`
				HasPrev     bool `json:"has_prev"`
			} `json:"pagination"`
		} `json:"data"`
	}

	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, "Success", response.Message)
	assert.NotEmpty(t, response.Data.Spaces, "Should return public spaces")

	for _, space := range response.Data.Spaces {
		assert.False(t, space.PrivacyStatus, "Only public spaces should be returned")
	}
}

func TestListApiKeys(t *testing.T) {
	router := setupSpaceRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/spaces/1/api-keys", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		Status  int                      `json:"status"`
		Message string                   `json:"message"`
		Data    []map[string]interface{} `json:"data"`
	}

	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, "Success", response.Message)
	assert.NotEmpty(t, response.Data, "Should return API keys")
}

func TestCreateApiKey(t *testing.T) {
	router := setupSpaceRouter()

	apiKeyData := struct {
		Name      string     `json:"name"`
		ExpiresAt *time.Time `json:"expires_at"`
	}{
		Name: "Test API Key",
	}

	payload, _ := json.Marshal(apiKeyData)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/spaces/1/api-keys", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response struct {
		Status  int                    `json:"status"`
		Message string                 `json:"message"`
		Data    map[string]interface{} `json:"data"`
	}

	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, "API key created successfully", response.Message)
	assert.Equal(t, "Test API Key", response.Data["name"])
	assert.Contains(t, response.Data["key"].(string), "sk-")
}

func TestGetApiKey(t *testing.T) {
	router := setupSpaceRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/spaces/1/api-keys/1", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		Status  int                    `json:"status"`
		Message string                 `json:"message"`
		Data    map[string]interface{} `json:"data"`
	}

	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, "Success", response.Message)
	assert.Equal(t, "1", response.Data["id"])
	assert.Equal(t, "Default API Key", response.Data["name"])
}

func TestDeleteApiKey(t *testing.T) {
	router := setupSpaceRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/spaces/1/api-keys/1", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
	}

	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, "API key deleted successfully", response.Message)
}

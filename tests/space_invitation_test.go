package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type SpaceInvitation struct {
	ID            uint      `json:"id"`
	SpaceID       uint      `json:"space_id"`
	InviterID     uint      `json:"inviter_id"`
	InvitedUserID uint      `json:"invited_user_id"`
	SpaceRoleID   uint      `json:"space_role_id"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

var mockSpaceInvitations = []SpaceInvitation{
	{
		ID:            1,
		SpaceID:       1,
		InviterID:     1,
		InvitedUserID: 2,
		SpaceRoleID:   2, // Editor
		Status:        "pending",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	},
	{
		ID:            2,
		SpaceID:       2,
		InviterID:     1,
		InvitedUserID: 3,
		SpaceRoleID:   3, // Viewer
		Status:        "pending",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	},
	{
		ID:            3,
		SpaceID:       1,
		InviterID:     1,
		InvitedUserID: 4,
		SpaceRoleID:   2, // Editor
		Status:        "accepted",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	},
}

func mock_Auth_Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	}
}

func pointer_ToUint(i uint) *uint {
	return &i
}

type Space_User struct {
	UserID      uint  `json:"user_id"`
	SpaceID     uint  `json:"space_id"`
	SpaceRoleID *uint `json:"space_role_id"`
}

var mock_Space_Users = []Space_User{
	{
		UserID:      1,
		SpaceID:     1,
		SpaceRoleID: pointer_ToUint(1), // Owner
	},
	{
		UserID:      2,
		SpaceID:     1,
		SpaceRoleID: pointer_ToUint(2), // Editor
	},
	{
		UserID:      1,
		SpaceID:     2,
		SpaceRoleID: pointer_ToUint(1), // Owner
	},
}

func setupSpaceInvitationRouter() *gin.Engine {
	r := gin.Default()

	r.Use(mock_Auth_Middleware())

	r.GET("/space-invitations", GetAllSpaceInvitationsHandler)
	r.GET("/space-invitations/:id", GetSpaceInvitationByIDHandler)
	r.GET("/space-invitations/count", GetInvitationCountHandler)
	r.PUT("/space-invitations/:id", UpdateSpaceInvitationHandler)
	r.PATCH("/space-invitations/:id", PatchSpaceInvitationHandler)
	r.DELETE("/space-invitations/:id", DeleteSpaceInvitationHandler)
	r.PUT("/space-invitations/:id/accept", AcceptInvitationHandler)
	r.PUT("/space-invitations/:id/reject", RejectInvitationHandler)

	return r
}

// GET /space-invitations - Retrieve all space invitations
func GetAllSpaceInvitationsHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Success",
		"data":    mockSpaceInvitations,
	})
}

// GET /space-invitations/:id - Retrieve one space invitation
func GetSpaceInvitationByIDHandler(c *gin.Context) {
	invitationID := c.Param("id")

	var foundInvitation *SpaceInvitation
	for _, invitation := range mockSpaceInvitations {
		if fmt.Sprintf("%d", invitation.ID) == invitationID {
			foundInvitation = &invitation
			break
		}
	}

	if foundInvitation == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "Invitation not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Success",
		"data":    foundInvitation,
	})
}

// GET /space-invitations/count - Get count of invitations for current user
func GetInvitationCountHandler(c *gin.Context) {
	userID, _ := c.Get("user_id")
	count := 0

	for _, invitation := range mockSpaceInvitations {
		if invitation.InvitedUserID == userID.(uint) && invitation.Status == "pending" {
			count++
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Invitation count retrieved successfully",
		"data": gin.H{
			"count": count,
		},
	})
}

// PUT /space-invitations/:id - Update a space invitation
func UpdateSpaceInvitationHandler(c *gin.Context) {
	invitationID := c.Param("id")
	var updatedInvitation SpaceInvitation

	if err := c.ShouldBindJSON(&updatedInvitation); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Invalid request",
		})
		return
	}

	found := false
	for i, invitation := range mockSpaceInvitations {
		if fmt.Sprintf("%d", invitation.ID) == invitationID {
			updatedInvitation.ID = invitation.ID
			updatedInvitation.CreatedAt = invitation.CreatedAt
			updatedInvitation.UpdatedAt = time.Now()
			mockSpaceInvitations[i] = updatedInvitation
			found = true
			break
		}
	}

	if !found {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "Invitation not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Updated",
		"data":    updatedInvitation,
	})
}

// PATCH /space-invitations/:id - Patch a space invitation
func PatchSpaceInvitationHandler(c *gin.Context) {
	invitationID := c.Param("id")
	var patchData map[string]interface{}

	if err := c.ShouldBindJSON(&patchData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Invalid request",
		})
		return
	}

	found := false
	var updatedInvitation SpaceInvitation
	for i, invitation := range mockSpaceInvitations {
		if fmt.Sprintf("%d", invitation.ID) == invitationID {
			updatedInvitation = invitation

			if status, ok := patchData["status"].(string); ok {
				updatedInvitation.Status = status
			}
			if spaceRoleID, ok := patchData["space_role_id"].(float64); ok {
				updatedInvitation.SpaceRoleID = uint(spaceRoleID)
			}

			updatedInvitation.UpdatedAt = time.Now()
			mockSpaceInvitations[i] = updatedInvitation
			found = true
			break
		}
	}

	if !found {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "Invitation not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Updated",
		"data":    updatedInvitation,
	})
}

// DELETE /space-invitations/:id - Delete a space invitation
func DeleteSpaceInvitationHandler(c *gin.Context) {
	invitationID := c.Param("id")

	found := false
	for i, invitation := range mockSpaceInvitations {
		if fmt.Sprintf("%d", invitation.ID) == invitationID {
			mockSpaceInvitations = append(mockSpaceInvitations[:i], mockSpaceInvitations[i+1:]...)
			found = true
			break
		}
	}

	if !found {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "Invitation not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Deleted",
	})
}

// PUT /space-invitations/:id/accept - Accept an invitation
func AcceptInvitationHandler(c *gin.Context) {
	invitationID := c.Param("id")
	userID, _ := c.Get("user_id")

	found := false
	var invitation SpaceInvitation
	for i, inv := range mockSpaceInvitations {
		if fmt.Sprintf("%d", inv.ID) == invitationID && inv.InvitedUserID == userID.(uint) {
			invitation = inv
			invitation.Status = "accepted"
			invitation.UpdatedAt = time.Now()
			mockSpaceInvitations[i] = invitation
			found = true
			break
		}
	}

	if !found {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "Invitation not found or not authorized",
		})
		return
	}

	// Add user to space members
	mock_Space_Users = append(mock_Space_Users, Space_User{
		UserID:      userID.(uint),
		SpaceID:     invitation.SpaceID,
		SpaceRoleID: pointer_ToUint(invitation.SpaceRoleID),
	})

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Invitation accepted successfully",
		"data": gin.H{
			"ok": "yes",
		},
	})
}

// PUT /space-invitations/:id/reject - Reject an invitation
func RejectInvitationHandler(c *gin.Context) {
	invitationID := c.Param("id")
	userID, _ := c.Get("user_id")

	found := false
	for i, invitation := range mockSpaceInvitations {
		if fmt.Sprintf("%d", invitation.ID) == invitationID && invitation.InvitedUserID == userID.(uint) {
			invitation.Status = "rejected"
			invitation.UpdatedAt = time.Now()
			mockSpaceInvitations[i] = invitation
			found = true
			break
		}
	}

	if !found {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "Invitation not found or not authorized",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Invitation rejected successfully",
		"data": gin.H{
			"ok": "yes",
		},
	})
}

// TestGetSpaceInvitations tests GET /space-invitations
func TestGetSpaceInvitations(t *testing.T) {
	router := setupSpaceInvitationRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/space-invitations", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		Status  int               `json:"status"`
		Message string            `json:"message"`
		Data    []SpaceInvitation `json:"data"`
	}

	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, "Success", response.Message)
	assert.Equal(t, len(mockSpaceInvitations), len(response.Data))
}

// TestGetSpaceInvitation tests GET /space-invitations/:id
func TestGetSpaceInvitation(t *testing.T) {
	router := setupSpaceInvitationRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/space-invitations/1", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		Status  int             `json:"status"`
		Message string          `json:"message"`
		Data    SpaceInvitation `json:"data"`
	}

	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, "Success", response.Message)
	assert.Equal(t, uint(1), response.Data.ID)
	assert.Equal(t, uint(1), response.Data.SpaceID)
}

// TestGetInvitationCount tests GET /space-invitations/count
func TestGetInvitationCount(t *testing.T) {
	router := setupSpaceInvitationRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/space-invitations/count", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
		Data    struct {
			Count int `json:"count"`
		} `json:"data"`
	}

	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, "Invitation count retrieved successfully", response.Message)
	assert.GreaterOrEqual(t, response.Data.Count, 0)
}

// TestUpdateSpaceInvitation tests PUT /space-invitations/:id
func TestUpdateSpaceInvitation(t *testing.T) {
	router := setupSpaceInvitationRouter()

	updatedInvitation := SpaceInvitation{
		SpaceID:       1,
		InviterID:     1,
		InvitedUserID: 2,
		SpaceRoleID:   3, // Change role to Viewer
		Status:        "pending",
	}

	payload, _ := json.Marshal(updatedInvitation)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/space-invitations/1", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		Status  int             `json:"status"`
		Message string          `json:"message"`
		Data    SpaceInvitation `json:"data"`
	}

	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, "Updated", response.Message)
	assert.Equal(t, uint(3), response.Data.SpaceRoleID)
}

// TestPatchSpaceInvitation tests PATCH /space-invitations/:id
func TestPatchSpaceInvitation(t *testing.T) {
	router := setupSpaceInvitationRouter()

	patchData := map[string]interface{}{
		"status": "cancelled",
	}

	payload, _ := json.Marshal(patchData)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/space-invitations/1", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		Status  int             `json:"status"`
		Message string          `json:"message"`
		Data    SpaceInvitation `json:"data"`
	}

	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, "Updated", response.Message)
	assert.Equal(t, "cancelled", response.Data.Status)
}

// TestDeleteSpaceInvitation tests DELETE /space-invitations/:id
func TestDeleteSpaceInvitation(t *testing.T) {
	router := setupSpaceInvitationRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/space-invitations/2", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
	}

	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, "Deleted", response.Message)
}

// TestAcceptInvitation tests PUT /space-invitations/:id/accept
func TestAcceptInvitation(t *testing.T) {
	mockSpaceInvitations = []SpaceInvitation{
		{
			ID:            1,
			SpaceID:       1,
			InviterID:     1,
			InvitedUserID: 1, // Match test user ID
			SpaceRoleID:   2,
			Status:        "pending",
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		},
	}

	router := setupSpaceInvitationRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/space-invitations/1/accept", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
		Data    struct {
			Ok string `json:"ok"`
		} `json:"data"`
	}

	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, "Invitation accepted successfully", response.Message)
	assert.Equal(t, "yes", response.Data.Ok)

	found := false
	for _, inv := range mockSpaceInvitations {
		if inv.ID == 1 {
			assert.Equal(t, "accepted", inv.Status)
			found = true
		}
	}
	assert.True(t, found)
}

// TestRejectInvitation tests PUT /space-invitations/:id/reject
func TestRejectInvitation(t *testing.T) {
	mockSpaceInvitations = []SpaceInvitation{
		{
			ID:            1,
			SpaceID:       1,
			InviterID:     1,
			InvitedUserID: 1, // Match test user ID
			SpaceRoleID:   2,
			Status:        "pending",
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		},
	}

	router := setupSpaceInvitationRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/space-invitations/1/reject", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
		Data    struct {
			Ok string `json:"ok"`
		} `json:"data"`
	}

	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, "Invitation rejected successfully", response.Message)
	assert.Equal(t, "yes", response.Data.Ok)

	found := false
	for _, inv := range mockSpaceInvitations {
		if inv.ID == 1 {
			assert.Equal(t, "rejected", inv.Status)
			found = true
		}
	}
	assert.True(t, found)
}

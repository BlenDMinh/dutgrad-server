package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type UserQuerySession struct {
	ID            uint          `json:"id"`
	UserID        *uint         `json:"user_id"`
	SpaceID       uint          `json:"space_id"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
	TempMessage   *string       `json:"temp_message"`
	ChatHistories []ChatHistory `json:"chat_histories"`
}

type ChatHistory struct {
	ID        uint            `json:"id"`
	SessionID uint            `json:"session_id"`
	Message   json.RawMessage `json:"message"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

type UserQuery struct {
	ID             uint      `json:"id"`
	QuerySessionID uint      `json:"query_session_id"`
	Query          string    `json:"query"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

var mockUserID uint = 1
var mockSpace = struct {
	ID          uint
	Name        string
	Description string
}{
	ID:          1,
	Name:        "Test Space",
	Description: "A test space for querying",
}

var mockQuerySessions = []UserQuerySession{
	{
		ID:        1,
		UserID:    &mockUserID,
		SpaceID:   mockSpace.ID,
		CreatedAt: time.Now().Add(-48 * time.Hour),
		UpdatedAt: time.Now().Add(-48 * time.Hour),
	},
	{
		ID:        2,
		UserID:    &mockUserID,
		SpaceID:   mockSpace.ID,
		CreatedAt: time.Now().Add(-24 * time.Hour),
		UpdatedAt: time.Now().Add(-24 * time.Hour),
	},
}

var mockChatHistories = []ChatHistory{
	{
		ID:        1,
		SessionID: 1,
		Message:   json.RawMessage(`{"type":"human","content":"How do I use this API?"}`),
		CreatedAt: time.Now().Add(-48 * time.Hour),
		UpdatedAt: time.Now().Add(-48 * time.Hour),
	},
	{
		ID:        2,
		SessionID: 1,
		Message:   json.RawMessage(`{"type":"ai","content":"Here's how you use the API..."}`),
		CreatedAt: time.Now().Add(-48 * time.Hour),
		UpdatedAt: time.Now().Add(-48 * time.Hour),
	},
	{
		ID:        3,
		SessionID: 2,
		Message:   json.RawMessage(`{"type":"human","content":"Thật không?"}`),
		CreatedAt: time.Now().Add(-24 * time.Hour),
		UpdatedAt: time.Now().Add(-24 * time.Hour),
	},
	{
		ID:        4,
		SessionID: 2,
		Message:   json.RawMessage(`{"type":"ai","content":"Chắc chắn rồi! 😊 Tôi đã xác nhận rằng không có hồ sơ nào về Nguyễn Thị Thu Hà trong lớp 21TCLC_Nhat1."}`),
		CreatedAt: time.Now().Add(-24 * time.Hour),
		UpdatedAt: time.Now().Add(-24 * time.Hour),
	},
}

var mockTempMessages = map[uint]*string{
	1: nil,
	2: func() *string { s := "This is a temporary message"; return &s }(),
}

var mockUserQueries = []UserQuery{
	{
		ID:             1,
		QuerySessionID: 1,
		Query:          "How do I use this API?",
		CreatedAt:      time.Now().Add(-48 * time.Hour),
		UpdatedAt:      time.Now().Add(-48 * time.Hour),
	},
	{
		ID:             2,
		QuerySessionID: 2,
		Query:          "What are the limitations?",
		CreatedAt:      time.Now().Add(-24 * time.Hour),
		UpdatedAt:      time.Now().Add(-24 * time.Hour)},
}

func mockAuthMiddle() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("user_id", mockUserID)
		c.Next()
	}
}

func setupQuerySessionRouter() *gin.Engine {
	r := gin.Default()
	r.Use(mockAuthMiddle())

	sessionGroup := r.Group("/user-query-sessions")
	{
		sessionGroup.POST("/begin-chat-session", BeginChatSessionHandler)
		sessionGroup.GET("/me", GetMyChatSessionsHandler)
		sessionGroup.HEAD("/me", CountMyChatSessionsHandler)
		sessionGroup.GET("/:id/temp-message", GetTempMessageByIDHandler)
		sessionGroup.GET("/:id/history", GetChatHistoryHandler)
		sessionGroup.DELETE("/:id/history", ClearChatHistoryHandler)
	}

	queryGroup := r.Group("/user-query")
	{
		queryGroup.POST("/ask", AskHandler)
	}

	return r
}

func BeginChatSessionHandler(c *gin.Context) {
	var req struct {
		SpaceID uint `json:"space_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Invalid request",
			"error":   err.Error(),
		})
		return
	}

	// Validate space_id
	if req.SpaceID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Invalid request",
			"error":   "space_id is required",
		})
		return
	}

	userID, _ := c.Get("user_id")
	uid := userID.(uint)

	// Create a new session
	newSession := UserQuerySession{
		ID:        uint(len(mockQuerySessions) + 1),
		UserID:    &uid,
		SpaceID:   req.SpaceID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Session created successfully",
		"data":    newSession,
	})
}

func GetMyChatSessionsHandler(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid := userID.(uint)

	var userSessions []UserQuerySession
	for _, session := range mockQuerySessions {
		if session.UserID != nil && *session.UserID == uid {
			userSessions = append(userSessions, session)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Fetched sessions successfully",
		"data":    userSessions,
	})
}

func CountMyChatSessionsHandler(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid := userID.(uint)

	var count int
	for _, session := range mockQuerySessions {
		if session.UserID != nil && *session.UserID == uid {
			count++
		}
	}

	c.Header("X-Total-Count", fmt.Sprintf("%d", count))
	c.Status(http.StatusOK)
}

func GetTempMessageByIDHandler(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Invalid session ID",
			"error":   err.Error(),
		})
		return
	}

	tempMessage, exists := mockTempMessages[uint(id)]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "Session not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Fetched temp message successfully",
		"data":    tempMessage,
	})
}

func GetChatHistoryHandler(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Invalid session ID",
			"error":   err.Error(),
		})
		return
	}

	userID, _ := c.Get("user_id")

	// Check if session exists and belongs to the user
	var sessionFound bool
	for _, session := range mockQuerySessions {
		if session.ID == uint(id) && session.UserID != nil && *session.UserID == userID.(uint) {
			sessionFound = true
			break
		}
	}

	if !sessionFound {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "Session not found or unauthorized",
		})
		return
	}

	var history []map[string]interface{}
	for _, chatHistory := range mockChatHistories {
		if chatHistory.SessionID == uint(id) {
			var msg map[string]interface{}
			_ = json.Unmarshal(chatHistory.Message, &msg)
			history = append(history, map[string]interface{}{
				"id":         chatHistory.ID,
				"session_id": chatHistory.SessionID,
				"message":    msg,
				"created_at": chatHistory.CreatedAt,
				"updated_at": chatHistory.UpdatedAt,
			})
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Fetched chat history successfully",
		"data":    history,
	})
}

func ClearChatHistoryHandler(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Invalid session ID",
			"error":   err.Error(),
		})
		return
	}

	userID, _ := c.Get("user_id")

	// Check if session exists and belongs to the user
	var sessionFound bool
	for _, session := range mockQuerySessions {
		if session.ID == uint(id) && session.UserID != nil && *session.UserID == userID.(uint) {
			sessionFound = true
			break
		}
	}

	if !sessionFound {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "Session not found or unauthorized",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Chat history and session cleared successfully",
	})
}

func AskHandler(c *gin.Context) {
	var req struct {
		QuerySessionID uint   `json:"query_session_id"`
		Query          string `json:"query"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Invalid request",
			"error":   err.Error(),
		})
		return
	}

	// Additional validation to ensure both QuerySessionID and Query are provided
	if req.QuerySessionID == 0 || req.Query == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Invalid request",
			"error":   "query_session_id and query are required fields",
		})
		return
	}

	// Verify the session exists
	var sessionFound bool
	var sessionSpaceID uint
	for _, session := range mockQuerySessions {
		if session.ID == req.QuerySessionID {
			sessionFound = true
			sessionSpaceID = session.SpaceID
			break
		}
	}

	if !sessionFound {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "Session not found",
		})
		return
	}
	// Mock a response from the RAG service
	answer := fmt.Sprintf("Answer to: %s (for space %d)", req.Query, sessionSpaceID)

	// Create new chat history entries
	userMessage := ChatHistory{
		ID:        uint(len(mockChatHistories) + 1),
		SessionID: req.QuerySessionID,
		Message:   json.RawMessage(fmt.Sprintf(`{"type":"human","content":"%s"}`, req.Query)),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	aiMessage := ChatHistory{
		ID:        uint(len(mockChatHistories) + 2),
		SessionID: req.QuerySessionID,
		Message:   json.RawMessage(fmt.Sprintf(`{"type":"ai","content":"%s"}`, answer)),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Add to mock chat histories
	mockChatHistories = append(mockChatHistories, userMessage, aiMessage)

	// Create a new query
	newQuery := UserQuery{
		ID:             uint(len(mockUserQueries) + 1),
		QuerySessionID: req.QuerySessionID,
		Query:          req.Query,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Answer retrieved successfully",
		"data": gin.H{
			"answer": answer,
			"query":  newQuery,
		},
	})
}

// Test cases
func TestBeginChatSession(t *testing.T) {
	router := setupQuerySessionRouter()

	tests := []struct {
		name         string
		requestBody  map[string]interface{}
		expectedCode int
		expectedMsg  string
	}{
		{
			name: "Create new chat session",
			requestBody: map[string]interface{}{
				"space_id": 1,
			},
			expectedCode: http.StatusOK,
			expectedMsg:  "Session created successfully",
		},
		{
			name:         "Missing space_id",
			requestBody:  map[string]interface{}{},
			expectedCode: http.StatusBadRequest,
			expectedMsg:  "Invalid request",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			jsonBody, _ := json.Marshal(tc.requestBody)
			req, _ := http.NewRequest("POST", "/user-query-sessions/begin-chat-session", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedCode, w.Code)

			var response map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &response)
			assert.Contains(t, response["message"], tc.expectedMsg)

			if tc.expectedCode == http.StatusOK {
				data, _ := response["data"].(map[string]interface{})
				assert.NotNil(t, data["id"])

				spaceIDFloat, _ := data["space_id"].(float64)
				if spaceID, ok := tc.requestBody["space_id"].(int); ok {
					assert.Equal(t, float64(spaceID), spaceIDFloat)
				} else if spaceID, ok := tc.requestBody["space_id"].(float64); ok {
					assert.Equal(t, spaceID, spaceIDFloat)
				}
			}
		})
	}
}

func TestGetMyChatSessions(t *testing.T) {
	router := setupQuerySessionRouter()

	req, _ := http.NewRequest("GET", "/user-query-sessions/me", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Contains(t, response["message"], "Fetched sessions successfully")

	data, _ := response["data"].([]interface{})
	assert.GreaterOrEqual(t, len(data), 2)
}

func TestCountMyChatSessions(t *testing.T) {
	router := setupQuerySessionRouter()

	req, _ := http.NewRequest("HEAD", "/user-query-sessions/me", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	count := w.Header().Get("X-Total-Count")
	assert.NotEmpty(t, count)

	countInt, err := strconv.Atoi(count)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, countInt, 2)
}

func TestGetTempMessage(t *testing.T) {
	router := setupQuerySessionRouter()

	tests := []struct {
		name         string
		sessionID    string
		expectedCode int
		expectedMsg  string
		checkData    bool
	}{
		{
			name:         "Get temp message for session 2",
			sessionID:    "2",
			expectedCode: http.StatusOK,
			expectedMsg:  "Fetched temp message successfully",
			checkData:    true,
		},
		{
			name:         "Session not found",
			sessionID:    "999",
			expectedCode: http.StatusNotFound,
			expectedMsg:  "Session not found",
			checkData:    false,
		},
		{
			name:         "Invalid session ID",
			sessionID:    "invalid",
			expectedCode: http.StatusBadRequest,
			expectedMsg:  "Invalid session ID",
			checkData:    false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/user-query-sessions/"+tc.sessionID+"/temp-message", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedCode, w.Code)

			var response map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &response)
			assert.Contains(t, response["message"], tc.expectedMsg)

			if tc.checkData {
				assert.NotNil(t, response["data"])
				if tc.sessionID == "2" {
					assert.Equal(t, "This is a temporary message", response["data"])
				}
			}
		})
	}
}

func TestGetChatHistory(t *testing.T) {
	router := setupQuerySessionRouter()

	tests := []struct {
		name         string
		sessionID    string
		expectedCode int
		expectedMsg  string
		expectedLen  int
	}{
		{
			name:         "Get chat history for session 1",
			sessionID:    "1",
			expectedCode: http.StatusOK,
			expectedMsg:  "Fetched chat history successfully",
			expectedLen:  2,
		},
		{
			name:         "Get chat history for session 2",
			sessionID:    "2",
			expectedCode: http.StatusOK,
			expectedMsg:  "Fetched chat history successfully",
			expectedLen:  2,
		},
		{
			name:         "Session not found",
			sessionID:    "999",
			expectedCode: http.StatusNotFound,
			expectedMsg:  "Session not found",
			expectedLen:  0,
		},
		{
			name:         "Invalid session ID",
			sessionID:    "invalid",
			expectedCode: http.StatusBadRequest,
			expectedMsg:  "Invalid session ID",
			expectedLen:  0,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/user-query-sessions/"+tc.sessionID+"/history", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedCode, w.Code)

			var response map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &response)
			assert.Contains(t, response["message"], tc.expectedMsg)

			if tc.expectedCode == http.StatusOK {
				data, _ := response["data"].([]interface{})
				assert.Equal(t, tc.expectedLen, len(data))

				// Verify message format
				if tc.sessionID == "2" {
					item1, _ := data[0].(map[string]interface{})
					message1, _ := item1["message"].(map[string]interface{})
					assert.Equal(t, "human", message1["type"])
					assert.Equal(t, "Thật không?", message1["content"])

					item2, _ := data[1].(map[string]interface{})
					message2, _ := item2["message"].(map[string]interface{})
					assert.Equal(t, "ai", message2["type"])
					assert.Contains(t, message2["content"], "Chắc chắn rồi!")
					assert.Contains(t, message2["content"], "Nguyễn Thị Thu Hà")
				}
			}
		})
	}
}

func TestClearChatHistory(t *testing.T) {
	router := setupQuerySessionRouter()

	tests := []struct {
		name         string
		sessionID    string
		expectedCode int
		expectedMsg  string
	}{
		{
			name:         "Clear chat history for session 1",
			sessionID:    "1",
			expectedCode: http.StatusOK,
			expectedMsg:  "Chat history and session cleared successfully",
		},
		{
			name:         "Session not found",
			sessionID:    "999",
			expectedCode: http.StatusNotFound,
			expectedMsg:  "Session not found",
		},
		{
			name:         "Invalid session ID",
			sessionID:    "invalid",
			expectedCode: http.StatusBadRequest,
			expectedMsg:  "Invalid session ID",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req, _ := http.NewRequest("DELETE", "/user-query-sessions/"+tc.sessionID+"/history", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedCode, w.Code)

			var response map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &response)
			assert.Contains(t, response["message"], tc.expectedMsg)
		})
	}
}

func TestAsk(t *testing.T) {
	router := setupQuerySessionRouter()
	tests := []struct {
		name         string
		requestBody  map[string]interface{}
		expectedCode int
		expectedMsg  string
	}{
		{
			name: "Ask a question with valid session",
			requestBody: map[string]interface{}{
				"query_session_id": 1,
				"query":            "How does this work?",
			},
			expectedCode: http.StatusOK,
			expectedMsg:  "Answer retrieved successfully",
		},
		{
			name: "Session not found",
			requestBody: map[string]interface{}{
				"query_session_id": 999,
				"query":            "How does this work?",
			},
			expectedCode: http.StatusNotFound,
			expectedMsg:  "Session not found",
		},
		{
			name: "Missing query",
			requestBody: map[string]interface{}{
				"query_session_id": 1,
			},
			expectedCode: http.StatusBadRequest,
			expectedMsg:  "Invalid request",
		},
		{
			name: "Missing session ID",
			requestBody: map[string]interface{}{
				"query": "How does this work?",
			},
			expectedCode: http.StatusBadRequest,
			expectedMsg:  "Invalid request",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			jsonBody, _ := json.Marshal(tc.requestBody)
			req, _ := http.NewRequest("POST", "/user-query/ask", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedCode, w.Code)

			var response map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &response)
			assert.Contains(t, response["message"], tc.expectedMsg)

			if tc.expectedCode == http.StatusOK {
				data, _ := response["data"].(map[string]interface{})
				assert.NotNil(t, data["answer"])
				assert.NotNil(t, data["query"])

				query, _ := data["query"].(map[string]interface{})
				assert.Equal(t, tc.requestBody["query"], query["query"])

				if sessionID, ok := tc.requestBody["query_session_id"].(int); ok {
					querySessionID, _ := query["query_session_id"].(float64)
					assert.Equal(t, float64(sessionID), querySessionID)
				} else if sessionID, ok := tc.requestBody["query_session_id"].(float64); ok {
					querySessionID, _ := query["query_session_id"].(float64)
					assert.Equal(t, sessionID, querySessionID)
				}

				if tc.name == "Ask a question in Vietnamese" {
					var sessionID float64
					if id, ok := tc.requestBody["query_session_id"].(int); ok {
						sessionID = float64(id)
					} else if id, ok := tc.requestBody["query_session_id"].(float64); ok {
						sessionID = id
					}

					historyReq, _ := http.NewRequest("GET", fmt.Sprintf("/user-query-sessions/%d/history",
						int(sessionID)), nil)
					historyW := httptest.NewRecorder()
					router.ServeHTTP(historyW, historyReq)

					var historyResponse map[string]interface{}
					json.Unmarshal(historyW.Body.Bytes(), &historyResponse)

					historyData, _ := historyResponse["data"].([]interface{})

					if len(historyData) >= 2 {
						lastMsg, _ := historyData[len(historyData)-1].(map[string]interface{})
						lastMsgContent, _ := lastMsg["message"].(map[string]interface{})
						assert.Equal(t, "ai", lastMsgContent["type"])

						secondLastMsg, _ := historyData[len(historyData)-2].(map[string]interface{})
						secondLastMsgContent, _ := secondLastMsg["message"].(map[string]interface{})
						assert.Equal(t, "human", secondLastMsgContent["type"])
						assert.Equal(t, tc.requestBody["query"], secondLastMsgContent["content"])
					}
				}
			}
		})
	}
}

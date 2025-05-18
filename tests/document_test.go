package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type Document struct {
	ID               uint   `json:"id"`
	SpaceID          uint   `json:"space_id"`
	Name             string `json:"name"`
	Description      string `json:"description"`
	MimeType         string `json:"mime_type"`
	Size             int64  `json:"size"`
	ProcessingStatus int    `json:"processing_status"`
	S3URL            string `json:"s3_url"`
	PrivacyStatus    bool   `json:"privacy_status"`
}

var mockDocuments = []Document{
	{
		ID:               1,
		SpaceID:          1,
		Name:             "document1",
		Description:      "Tài liệu mẫu 1",
		MimeType:         "application/pdf",
		Size:             1024,
		ProcessingStatus: 1,
		S3URL:            "https://s3-bucket.com/document1.pdf",
		PrivacyStatus:    true,
	},
	{
		ID:               2,
		SpaceID:          1,
		Name:             "document2",
		Description:      "Tài liệu mẫu 2",
		MimeType:         "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
		Size:             2048,
		ProcessingStatus: 1,
		S3URL:            "https://s3-bucket.com/document2.sheet",
		PrivacyStatus:    true,
	},
	{
		ID:               3,
		SpaceID:          2,
		Name:             "document3",
		Description:      "Tài liệu mẫu 3",
		MimeType:         "application/pdf",
		Size:             3072,
		ProcessingStatus: 1,
		S3URL:            "https://s3-bucket.com/document3.pdf",
		PrivacyStatus:    true,
	},
}

var userRoles = map[uint]map[uint]string{
	1: {1: "owner", 2: "viewer"},
	2: {1: "editor", 2: "owner"},
}

func setupDocumentRouter() *gin.Engine {
	r := gin.Default()

	documents := r.Group("/documents")
	{
		documents.GET("", GetDocumentsHandler)
		documents.GET("/:id", GetDocumentHandler)
		documents.PUT("/:id", UpdateDocumentHandler)
		documents.PATCH("/:id", PatchDocumentHandler)
		documents.POST("/upload", UploadDocumentHandler)
		documents.DELETE("/:id", DeleteDocumentHandler)
	}

	spaces := r.Group("/space")
	{
		spaces.GET("/:space_id/documents", GetDocumentsBySpaceIDHandler)
	}

	return r
}

func GetDocumentsHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Success",
		"data": gin.H{
			"documents": mockDocuments,
		},
	})
}

func GetDocumentHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Invalid document ID",
		})
		return
	}

	for _, doc := range mockDocuments {
		if int(doc.ID) == id {
			c.JSON(http.StatusOK, gin.H{
				"status":  http.StatusOK,
				"message": "Success",
				"data": gin.H{
					"document": doc,
				},
			})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{
		"status":  http.StatusNotFound,
		"message": "Document not found",
	})
}

func UpdateDocumentHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Invalid document ID",
		})
		return
	}

	var updateData map[string]interface{}
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Invalid request body",
		})
		return
	}

	for i, doc := range mockDocuments {
		if int(doc.ID) == id {
			c.JSON(http.StatusOK, gin.H{
				"status":  http.StatusOK,
				"message": "Document updated successfully",
				"data": gin.H{
					"document": mockDocuments[i],
				},
			})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{
		"status":  http.StatusNotFound,
		"message": "Document not found",
	})
}

func PatchDocumentHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Invalid document ID",
		})
		return
	}

	var patchData map[string]interface{}
	if err := c.ShouldBindJSON(&patchData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Invalid request body",
		})
		return
	}

	for i, doc := range mockDocuments {
		if int(doc.ID) == id {
			c.JSON(http.StatusOK, gin.H{
				"status":  http.StatusOK,
				"message": "Document patched successfully",
				"data": gin.H{
					"document": mockDocuments[i],
				},
			})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{
		"status":  http.StatusNotFound,
		"message": "Document not found",
	})
}

func UploadDocumentHandler(c *gin.Context) {
	userID := uint(1)
	spaceIDStr := c.PostForm("space_id")
	spaceID, err := strconv.Atoi(spaceIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Invalid space ID",
		})
		return
	}

	role, exists := userRoles[userID][uint(spaceID)]
	if !exists || (role != "owner" && role != "editor") {
		c.JSON(http.StatusForbidden, gin.H{
			"status":  http.StatusForbidden,
			"message": "You don't have permission to upload documents to this space",
		})
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "No file uploaded",
		})
		return
	}
	defer file.Close()

	fileName := header.Filename
	fileSize := header.Size
	mimeType := c.Request.Header.Get("Mime-Type")
	description := c.PostForm("description")

	if fileSize > 10*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "File too large",
		})
		return
	}

	newDocument := Document{
		ID:               uint(len(mockDocuments) + 1),
		SpaceID:          uint(spaceID),
		Name:             fileName,
		Description:      description,
		MimeType:         mimeType,
		Size:             fileSize,
		ProcessingStatus: 0,
		S3URL:            fmt.Sprintf("https://s3-bucket.com/%s", fileName),
		PrivacyStatus:    true,
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "File uploaded successfully",
		"data": gin.H{
			"document": newDocument,
		},
	})
}

func DeleteDocumentHandler(c *gin.Context) {
	userID := uint(1)

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Invalid document ID",
		})
		return
	}

	var targetDoc *Document
	for _, doc := range mockDocuments {
		if int(doc.ID) == id {
			targetDoc = &doc
			break
		}
	}

	if targetDoc == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "Document not found",
		})
		return
	}

	role, exists := userRoles[userID][targetDoc.SpaceID]
	if !exists || (role != "owner" && role != "editor") {
		c.JSON(http.StatusForbidden, gin.H{
			"status":  http.StatusForbidden,
			"message": "You don't have permission to delete this document",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Document deleted successfully",
	})
}

func GetDocumentsBySpaceIDHandler(c *gin.Context) {
	spaceIDStr := c.Param("space_id")
	spaceID, err := strconv.Atoi(spaceIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Invalid space ID",
		})
		return
	}

	// Lọc documents theo space_id
	var docsInSpace []Document
	for _, doc := range mockDocuments {
		if int(doc.SpaceID) == spaceID {
			docsInSpace = append(docsInSpace, doc)
		}
	}

	if docsInSpace == nil {
		docsInSpace = []Document{}
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Success",
		"data": gin.H{
			"documents": docsInSpace,
		},
	})
}

func createUploadRequest(t *testing.T, url, filePath, spaceID, description string) (*http.Request, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	writer.WriteField("space_id", spaceID)
	writer.WriteField("description", description)

	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return nil, err
	}
	io.Copy(part, file)

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Mime-Type", "application/pdf")

	return req, nil
}

// Test cases
func TestGetAllDocuments(t *testing.T) {
	router := setupDocumentRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/documents", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, float64(http.StatusOK), response["status"])
	assert.Equal(t, "Success", response["message"])

	data, exists := response["data"].(map[string]interface{})
	assert.True(t, exists)

	documents, exists := data["documents"].([]interface{})
	assert.True(t, exists)
	assert.Len(t, documents, len(mockDocuments))
}

func TestGetDocumentByID(t *testing.T) {
	router := setupDocumentRouter()

	testCases := []struct {
		name         string
		documentID   string
		expectedCode int
		checkData    bool
	}{
		{
			name:         "✅ Lấy tài liệu thành công",
			documentID:   "1",
			expectedCode: http.StatusOK,
			checkData:    true,
		},
		{
			name:         "❌ ID tài liệu không hợp lệ",
			documentID:   "invalid",
			expectedCode: http.StatusBadRequest,
			checkData:    false,
		},
		{
			name:         "❌ Tài liệu không tồn tại",
			documentID:   "999",
			expectedCode: http.StatusNotFound,
			checkData:    false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/documents/"+tc.documentID, nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedCode, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			if tc.checkData {
				assert.Equal(t, "Success", response["message"])
				data, exists := response["data"].(map[string]interface{})
				assert.True(t, exists)
				document, exists := data["document"].(map[string]interface{})
				assert.True(t, exists)
				assert.Equal(t, float64(1), document["id"])
			}
		})
	}
}

func TestUpdateDocument(t *testing.T) {
	router := setupDocumentRouter()

	updateBody := map[string]interface{}{
		"description":    "Mô tả cập nhật",
		"privacy_status": false,
	}
	jsonBody, _ := json.Marshal(updateBody)

	testCases := []struct {
		name         string
		documentID   string
		requestBody  *bytes.Buffer
		expectedCode int
	}{
		{
			name:         "✅ Cập nhật tài liệu thành công",
			documentID:   "1",
			requestBody:  bytes.NewBuffer(jsonBody),
			expectedCode: http.StatusOK,
		},
		{
			name:         "❌ Dữ liệu không hợp lệ",
			documentID:   "1",
			requestBody:  bytes.NewBuffer([]byte("Invalid JSON")),
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "❌ Tài liệu không tồn tại",
			documentID:   "999",
			requestBody:  bytes.NewBuffer(jsonBody),
			expectedCode: http.StatusNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("PUT", "/documents/"+tc.documentID, tc.requestBody)
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedCode, w.Code)
		})
	}
}

func TestDeleteDocument(t *testing.T) {
	router := setupDocumentRouter()

	testCases := []struct {
		name         string
		documentID   string
		expectedCode int
	}{
		{
			name:         "✅ Xóa tài liệu thành công",
			documentID:   "1",
			expectedCode: http.StatusOK,
		},
		{
			name:         "❌ ID tài liệu không hợp lệ",
			documentID:   "invalid",
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "❌ Tài liệu không tồn tại",
			documentID:   "999",
			expectedCode: http.StatusNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("DELETE", "/documents/"+tc.documentID, nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedCode, w.Code)
		})
	}
}

func TestGetDocumentsBySpaceID(t *testing.T) {
	router := setupDocumentRouter()

	testCases := []struct {
		name         string
		spaceID      string
		expectedCode int
		expectedDocs int
	}{
		{
			name:         "✅ Lấy tài liệu trong space thành công",
			spaceID:      "1",
			expectedCode: http.StatusOK,
			expectedDocs: 2, // Space ID 1 có 2 documents
		},
		{
			name:         "✅ Space không có tài liệu",
			spaceID:      "3",
			expectedCode: http.StatusOK,
			expectedDocs: 0,
		},
		{
			name:         "❌ ID space không hợp lệ",
			spaceID:      "invalid",
			expectedCode: http.StatusBadRequest,
			expectedDocs: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/space/"+tc.spaceID+"/documents", nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedCode, w.Code)

			if tc.expectedCode == http.StatusOK {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)

				data, exists := response["data"].(map[string]interface{})
				assert.True(t, exists)

				documents, exists := data["documents"].([]interface{})
				assert.True(t, exists)
				assert.Len(t, documents, tc.expectedDocs)
			}
		})
	}
}

func TestUploadDocument(t *testing.T) {
	router := setupDocumentRouter()

	t.Run("✅ Upload tài liệu thành công", func(t *testing.T) {
		tempFile, err := os.CreateTemp("", "test-upload-*.pdf")
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(tempFile.Name())

		if _, err := tempFile.Write([]byte("Test document content")); err != nil {
			t.Fatal(err)
		}
		tempFile.Close()

		req, err := createUploadRequest(t, "/documents/upload", tempFile.Name(), "1", "Tài liệu test")
		if err != nil {
			t.Fatal(err)
		}

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Equal(t, "File uploaded successfully", response["message"])

		data, exists := response["data"].(map[string]interface{})
		assert.True(t, exists)

		document, exists := data["document"].(map[string]interface{})
		assert.True(t, exists)
		assert.Equal(t, float64(1), document["space_id"])
	})

	t.Run("❌ Space ID không hợp lệ", func(t *testing.T) {
		// Mock request không có space_id
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		writer.Close()

		req, _ := http.NewRequest("POST", "/documents/upload", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

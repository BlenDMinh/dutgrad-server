package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/BlenDMinh/dutgrad-server/databases/entities"
	"github.com/BlenDMinh/dutgrad-server/models/dtos"
	"github.com/BlenDMinh/dutgrad-server/services"
	"github.com/gin-gonic/gin"
)

type DocumentController struct {
	CrudController[entities.Document, uint]
	service      services.DocumentService
	spaceService services.SpaceService
}

func NewDocumentController(
	service services.DocumentService,
	spaceService services.SpaceService,
) *DocumentController {
	crudController := NewCrudController(service)
	return &DocumentController{
		CrudController: *crudController,
		service:        service,
		spaceService:   spaceService,
	}
}

func (c *DocumentController) GetBySpaceID(ctx *gin.Context) {
	spaceID, ok := ExtractID(ctx, "id")
	if !ok {
		return
	}

	documents, err := c.service.GetDocumentsBySpaceID(spaceID)
	if err != nil {
		HandleError(ctx, http.StatusInternalServerError, "Failed to retrieve documents", err)
		return
	}

	HandleSuccess(ctx, "Documents retrieved successfully", gin.H{"documents": documents})
}

func (c *DocumentController) UploadDocument(ctx *gin.Context) {
	userID, ok := ExtractID(ctx, "user_id")
	if !ok {
		return
	}

	var req dtos.DocumentUploadRequest
	spaceIDStr := ctx.Request.FormValue("space_id")
	spaceID, err := strconv.ParseUint(spaceIDStr, 10, 32)
	if err != nil {
		HandleError(ctx, http.StatusBadRequest, "Invalid space_id", err)
		return
	}
	req.SpaceID = uint(spaceID)
	req.Description = ctx.Request.FormValue("description")
	role, err := c.spaceService.GetUserRole(userID, req.SpaceID)
	if err != nil {
		HandleError(ctx, http.StatusInternalServerError, "Failed to get user role", err)
		return
	}

	if !role.IsOwner() && !role.IsEditor() {
		HandleError(ctx, http.StatusForbidden, "You are not allowed to import documents to this space", nil)
		return
	}

	file, err := ctx.FormFile("file")
	if err != nil {
		HandleError(ctx, http.StatusBadRequest, "Failed to get file", err)
		return
	}
	mimeType := ctx.Request.Header.Get("Mime-Type")

	document, err := c.service.UploadDocument(file, req.SpaceID, mimeType, req.Description)
	if err != nil {
		statusCode := http.StatusInternalServerError

		if strings.Contains(err.Error(), "document limit reached") ||
			strings.Contains(err.Error(), "file size exceeds the limit") {
			statusCode = http.StatusTooManyRequests
		}

		HandleError(ctx, statusCode, "Failed to upload document", err)
		return
	}

	HandleSuccess(ctx, "Document uploaded successfully", gin.H{"document": document})
}

func (c *DocumentController) GetUserDocumentCount(ctx *gin.Context) {
	userID, ok := ExtractID(ctx, "user_id")
	if !ok {
		return
	}

	count, err := c.service.CountUserDocuments(userID)
	if err != nil {
		HandleError(ctx, http.StatusInternalServerError, "Failed to count documents", err)
		return
	}

	ctx.Header("X-Document-Count", fmt.Sprintf("%d", count))

	HandleSuccess(ctx, "Document count retrieved successfully", gin.H{"count": count})
}

func (c *DocumentController) DeleteDocument(ctx *gin.Context) {
	userID, ok := ExtractID(ctx, "user_id")
	if !ok {
		return
	}

	docID, ok := ExtractID(ctx, "id")
	if !ok {
		return
	}

	document, err := c.service.GetById(docID)
	if err != nil {
		HandleError(ctx, http.StatusNotFound, "Document not found", err)
		return
	}
	role, err := c.spaceService.GetUserRole(userID, document.SpaceID)
	if err != nil {
		HandleError(ctx, http.StatusInternalServerError, "Failed to get user role", err)
		return
	}

	if !role.IsOwner() && !role.IsEditor() {
		HandleError(ctx, http.StatusForbidden, "You are not allowed to delete this document", nil)
		return
	}

	err = c.service.DeleteDocument(docID)
	if err != nil {
		HandleError(ctx, http.StatusInternalServerError, "Failed to delete document", err)
		return
	}

	HandleSuccess(ctx, "Document deleted successfully", nil)
}

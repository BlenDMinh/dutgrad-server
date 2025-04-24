package controllers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/BlenDMinh/dutgrad-server/databases/entities"
	"github.com/BlenDMinh/dutgrad-server/models"
	"github.com/BlenDMinh/dutgrad-server/services"
	"github.com/gin-gonic/gin"
)

type DocumentController struct {
	CrudController[entities.Document, uint]
	service *services.DocumentService
}

func NewDocumentController() *DocumentController {
	service := services.NewDocumentService()
	return &DocumentController{
		CrudController: *NewCrudController(service),
		service:        service,
	}
}

func (c *DocumentController) GetBySpaceID(ctx *gin.Context) {
	spaceIDStr := ctx.Param("space_id")
	spaceID, err := strconv.ParseUint(spaceIDStr, 10, 32)
	if err != nil {
		errMsg := err.Error()
		ctx.JSON(
			http.StatusInternalServerError,
			models.NewErrorResponse(
				http.StatusInternalServerError,
				"Invalid space_id",
				&errMsg,
			),
		)
		return
	}

	documents, err := c.service.GetDocumentsBySpaceID(uint(spaceID))
	if err != nil {
		errMsg := err.Error()
		ctx.JSON(
			http.StatusInternalServerError,
			models.NewErrorResponse(
				http.StatusInternalServerError,
				"Failed to retrieve documents",
				&errMsg,
			),
		)
		return
	}

	ctx.JSON(http.StatusOK, models.NewSuccessResponse(
		http.StatusOK,
		"Success",
		gin.H{"documents": documents},
	))
}

func (c *DocumentController) UploadDocument(ctx *gin.Context) {
	spaceIDStr := ctx.Request.FormValue("space_id")
	spaceID, err := strconv.ParseUint(spaceIDStr, 10, 32)
	if err != nil {
		errMsg := err.Error()
		ctx.JSON(
			http.StatusInternalServerError,
			models.NewErrorResponse(
				http.StatusBadRequest,
				"Invalid space_id",
				&errMsg,
			),
		)
		return
	}

	file, err := ctx.FormFile("file")
	if err != nil {
		errMsg := err.Error()
		ctx.JSON(
			http.StatusInternalServerError,
			models.NewErrorResponse(
				http.StatusBadRequest,
				"Failed to get file",
				&errMsg,
			),
		)
		return
	}

	document, err := c.service.UploadDocument(file, uint(spaceID))
	if err != nil {
		errMsg := err.Error()
		statusCode := http.StatusInternalServerError

		if strings.Contains(errMsg, "document limit reached") || strings.Contains(errMsg, "file size exceeds the limit") {
			statusCode = http.StatusTooManyRequests
		}

		ctx.JSON(
			statusCode,
			models.NewErrorResponse(
				statusCode,
				"Failed to upload document",
				&errMsg,
			),
		)
		return
	}

	ctx.JSON(http.StatusOK, models.NewSuccessResponse(
		http.StatusOK,
		"Success",
		gin.H{"document": document},
	))
}

func (c *DocumentController) DeleteDocument(ctx *gin.Context) {
	userIDInterface, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, models.NewErrorResponse(http.StatusUnauthorized, "Unauthorized", nil))
		return
	}
	userID := userIDInterface.(uint)

	docIDStr := ctx.Param("id")
	docID, err := strconv.ParseUint(docIDStr, 10, 32)
	if err != nil {
		errMsg := err.Error()
		ctx.JSON(http.StatusBadRequest, models.NewErrorResponse(http.StatusBadRequest, "Invalid document ID", &errMsg))
		return
	}

	document, err := c.service.GetById(uint(docID))
	if err != nil {
		errMsg := err.Error()
		ctx.JSON(http.StatusNotFound, models.NewErrorResponse(http.StatusNotFound, "Document not found", &errMsg))
		return
	}

	role, err := c.service.GetUserRoleInSpace(userID, document.SpaceID)
	if err != nil {
		errMsg := err.Error()
		ctx.JSON(http.StatusInternalServerError, models.NewErrorResponse(http.StatusInternalServerError, "Failed to get user role", &errMsg))
		return
	}

	if role != "owner" && role != "editor" {
		ctx.JSON(http.StatusForbidden, models.NewErrorResponse(http.StatusForbidden, "You are not allowed to delete this document", nil))
		return
	}

	err = c.service.DeleteDocument(uint(docID))
	if err != nil {
		errMsg := err.Error()
		ctx.JSON(http.StatusInternalServerError, models.NewErrorResponse(http.StatusInternalServerError, "Failed to delete document", &errMsg))
		return
	}

	ctx.JSON(http.StatusOK, models.NewSuccessResponse(http.StatusOK, "Document deleted successfully", gin.H{}))
}

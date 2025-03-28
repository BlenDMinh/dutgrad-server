package controllers

import (
	"net/http"
	"strconv"

	"github.com/BlenDMinh/dutgrad-server/databases/entities"
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
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid space_id"})
		return
	}

	documents, err := c.service.GetDocumentsBySpaceID(uint(spaceID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve documents"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"documents": documents})
}

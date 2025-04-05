package controllers

import (
	"net/http"
	"strconv"

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

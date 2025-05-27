package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/BlenDMinh/dutgrad-server/databases/entities"
	"github.com/BlenDMinh/dutgrad-server/models/dtos"
	"github.com/BlenDMinh/dutgrad-server/services"
	"github.com/gin-gonic/gin"
)

type CrudController[T entities.Entity, ID any] struct {
	service services.ICrudService[T, ID]
}

func NewCrudController[T entities.Entity, ID any](service services.ICrudService[T, ID]) *CrudController[T, ID] {
	return &CrudController[T, ID]{
		service: service,
	}
}

func (c *CrudController[T, ID]) getModel() *T {
	return new(T)
}

func (c *CrudController[T, ID]) parseID(ctx *gin.Context) (ID, error) {
	idStr := ctx.Param("id")
	var id ID
	var err error

	switch any(id).(type) {
	case int:
		var idInt int
		idInt, err = strconv.Atoi(idStr)
		id = any(idInt).(ID)
	case uint:
		var idUint uint64
		idUint, err = strconv.ParseUint(idStr, 10, 32)
		id = any(uint(idUint)).(ID)
	case uint8:
		var idUint8 uint64
		idUint8, err = strconv.ParseUint(idStr, 10, 8)
		id = any(uint8(idUint8)).(ID)
	case uint16:
		var idUint16 uint64
		idUint16, err = strconv.ParseUint(idStr, 10, 16)
		id = any(uint16(idUint16)).(ID)
	case uint32:
		var idUint32 uint64
		idUint32, err = strconv.ParseUint(idStr, 10, 32)
		id = any(uint32(idUint32)).(ID)
	case string:
		id = any(idStr).(ID)
	default:
		err = fmt.Errorf("unsupported ID type")
	}

	return id, err
}

func (c *CrudController[T, ID]) Retrieve(ctx *gin.Context) {
	pageStr := ctx.DefaultQuery("page", "1")
	pageSizeStr := ctx.DefaultQuery("page-size", "20")

	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)

	entities, err := c.service.GetAll(page, pageSize)
	if err != nil {
		HandleError(ctx, http.StatusInternalServerError, "Failed to retrieve entities", err)
		return
	}

	total, err := c.service.Count()
	if err != nil {
		HandleError(ctx, http.StatusInternalServerError, "Failed to count entities", err)
		return
	}

	pagination := dtos.PaginationResponse{
		CurrentPage: page,
		PageSize:    pageSize,
		TotalPages:  (total + int64(pageSize) - 1) / int64(pageSize),
		TotalItems:  total,
		HasNext:     int64(page*pageSize) < total,
		HasPrev:     page > 1,
	}

	HandleSuccess(ctx, "Entities retrieved successfully", gin.H{
		"data":       entities,
		"pagination": pagination,
	})
}

func (c *CrudController[T, ID]) RetrieveOne(ctx *gin.Context) {
	id, err := c.parseID(ctx)
	if err != nil {
		HandleError(ctx, http.StatusBadRequest, "Invalid ID format", err)
		return
	}

	entity, err := c.service.GetById(id)
	if err != nil {
		HandleError(ctx, http.StatusInternalServerError, "Failed to retrieve entity", err)
		return
	}

	HandleSuccess(ctx, "Entity retrieved successfully", entity)
}

func (c *CrudController[T, ID]) Create(ctx *gin.Context) {
	model := c.getModel()
	if !HandleBindJSON(ctx, model) {
		return
	}

	createdModel, err := c.service.Create(model)
	if err != nil {
		HandleError(ctx, http.StatusInternalServerError, "Failed to create entity", err)
		return
	}

	HandleCreated(ctx, "Entity created successfully", createdModel)
}

func (c *CrudController[T, ID]) Update(ctx *gin.Context) {
	id, err := c.parseID(ctx)
	if err != nil {
		HandleError(ctx, http.StatusBadRequest, "Invalid ID format", err)
		return
	}

	model := c.getModel()
	if !HandleBindJSON(ctx, model) {
		return
	}

	updatedModel, err := c.service.UpdateByID(id, model)
	if err != nil {
		HandleError(ctx, http.StatusInternalServerError, "Failed to update entity", err)
		return
	}

	HandleSuccess(ctx, "Entity updated successfully", updatedModel)
}

func (c *CrudController[T, ID]) Patch(ctx *gin.Context) {
	id, err := c.parseID(ctx)
	if err != nil {
		HandleError(ctx, http.StatusBadRequest, "Invalid ID format", err)
		return
	}

	model := c.getModel()
	if !HandleBindJSON(ctx, model) {
		return
	}

	patchedModel, err := c.service.PatchByID(id, model)
	if err != nil {
		HandleError(ctx, http.StatusInternalServerError, "Failed to patch entity", err)
		return
	}

	HandleSuccess(ctx, "Entity patched successfully", patchedModel)
}

func (c *CrudController[T, ID]) Delete(ctx *gin.Context) {
	id, err := c.parseID(ctx)
	if err != nil {
		HandleError(ctx, http.StatusBadRequest, "Invalid ID format", err)
		return
	}

	err = c.service.Delete(id)
	if err != nil {
		HandleError(ctx, http.StatusInternalServerError, "Failed to delete entity", err)
		return
	}

	HandleSuccess(ctx, "Entity deleted successfully", nil)
}

func (c *CrudController[T, ID]) RegisterCRUD(router gin.IRouter) {
	router.GET("", c.Retrieve)
	router.GET("/:id", c.RetrieveOne)
	router.POST("", c.Create)
	router.PUT("/:id", c.Update)
	router.PATCH("/:id", c.Patch)
	router.DELETE("/:id", c.Delete)
}

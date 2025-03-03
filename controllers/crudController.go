package controllers

import (
	"github.com/BlenDMinh/dutgrad-server/databases"
	"github.com/BlenDMinh/dutgrad-server/models"
	"github.com/gin-gonic/gin"
)

type CrudController[T any] struct{}

func (c *CrudController[T]) getModel() *T {
	return new(T)
}

func (c *CrudController[T]) Retrieve(ctx *gin.Context) {
	db := databases.GetDB()
	entities := []T{}
	dbctx := db.Find(&entities)

	if dbctx.Error != nil {
		errMsg := dbctx.Error.Error()
		ctx.JSON(500, models.NewErrorResponse(500, "Internal Server Error", &errMsg))
		return
	}

	ctx.JSON(200, models.NewSuccessResponse(200, "Retrieve", entities))
}

func (c *CrudController[T]) RetrieveOne(ctx *gin.Context) {
	db := databases.GetDB()
	entity := new(T)
	dbctx := db.First(entity, ctx.Param("id"))
	if dbctx.Error != nil {
		errMsg := dbctx.Error.Error()
		ctx.JSON(500, models.NewErrorResponse(500, "Internal Server Error", &errMsg))
		return
	}

	ctx.JSON(200, models.NewSuccessResponse(200, "Retrieve", entity))
}

func (c *CrudController[T]) Create(ctx *gin.Context) {
	model := c.getModel()
	if err := ctx.ShouldBindJSON(model); err != nil {
		errMsg := err.Error()
		ctx.JSON(400, models.NewErrorResponse(400, "Bad Request", &errMsg))
		return
	}

	db := databases.GetDB()
	dbctx := db.Create(model)
	if dbctx.Error != nil {
		errMsg := dbctx.Error.Error()
		ctx.JSON(500, models.NewErrorResponse(500, "Internal Server Error", &errMsg))
		return
	}

	ctx.JSON(201, models.NewSuccessResponse(201, "Created", model))
}

func (c *CrudController[T]) Update(ctx *gin.Context) {
	model := c.getModel()
	if err := ctx.ShouldBindJSON(model); err != nil {
		errMsg := err.Error()
		ctx.JSON(400, models.NewErrorResponse(400, "Bad Request", &errMsg))
		return
	}

	db := databases.GetDB()
	dbctx := db.Save(model)
	if dbctx.Error != nil {
		errMsg := dbctx.Error.Error()
		ctx.JSON(500, models.NewErrorResponse(500, "Internal Server Error", &errMsg))
		return
	}

	ctx.JSON(200, models.NewSuccessResponse(200, "Updated", model))
}

func (c *CrudController[T]) Delete(ctx *gin.Context) {
	model := c.getModel()
	db := databases.GetDB()
	dbctx := db.Delete(model, ctx.Param("id"))
	if dbctx.Error != nil {
		errMsg := dbctx.Error.Error()
		ctx.JSON(500, models.NewErrorResponse(500, "Internal Server Error", &errMsg))
		return
	}

	ctx.JSON(200, models.NewSuccessResponse(200, "Deleted", nil))
}

func (c *CrudController[T]) RegisterCRUD(router gin.IRouter) {
	router.GET("", c.Retrieve)
	router.GET("/:id", c.RetrieveOne)
	router.POST("", c.Create)
	router.PUT("/:id", c.Update)
	router.DELETE("/:id", c.Delete)
}

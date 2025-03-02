package controller

import (
	"github.com/BlenDMinh/dutgrad-server/database"
	"github.com/BlenDMinh/dutgrad-server/model"
	"github.com/gin-gonic/gin"
)


type CrudController[T any] struct {}

func (c *CrudController[T]) getModel() *T {
	return new(T)
}

func (c *CrudController[T]) Retrieve(ctx *gin.Context) {
	db := database.GetDB()
	entities := []T{}
	dbctx := db.Find(&entities)

	if dbctx.Error != nil {
		errMsg := dbctx.Error.Error()
		ctx.JSON(500, model.NewErrorResponse(500, "Internal Server Error", &errMsg))
	}

	ctx.JSON(200, model.NewSuccessResponse(200, "Retrieve", entities))
}

func (c *CrudController[T]) RetrieveOne(ctx *gin.Context) {
	db := database.GetDB()
	entity := new(T)
	dbctx := db.First(entity, ctx.Param("id"))
	if dbctx.Error != nil {
		errMsg := dbctx.Error.Error()
		ctx.JSON(500, model.NewErrorResponse(500, "Internal Server Error", &errMsg))
	}

	ctx.JSON(200, model.NewSuccessResponse(200, "Retrieve", entity))
}

func (c *CrudController[T]) Create(ctx *gin.Context) {
	model := c.getModel()
	ctx.JSON(200, gin.H{
		"message": "Create",
		"model":   model,
	})
}

func (c *CrudController[T]) Update(ctx *gin.Context) {
	model := c.getModel()
	ctx.JSON(200, gin.H{
		"message": "Update",
		"model":   model,
	})
}

func (c *CrudController[T]) Delete(ctx *gin.Context) {
	model := c.getModel()
	ctx.JSON(200, gin.H{
		"message": "Delete",
		"model":   model,
	})
}

func (c *CrudController[T]) Register(router gin.IRouter) {
	router.GET("", c.Retrieve)
	router.GET("/:id", c.RetrieveOne)
	router.POST("", c.Create)
	router.PUT("/:id", c.Update)
	router.DELETE("/:id", c.Delete)
}

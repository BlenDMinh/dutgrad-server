package controller

import "github.com/gin-gonic/gin"

type ICrudController interface {
	getModel() interface{}
}

type CrudController struct {
	ICrudController
}

func (c *CrudController) Retrieve(ctx *gin.Context) {
	model := c.getModel()
	ctx.JSON(200, gin.H{
		"message": "Retrieve",
		"model":   model,
	})
}

func (c *CrudController) Create(ctx *gin.Context) {
	model := c.getModel()
	ctx.JSON(200, gin.H{
		"message": "Create",
		"model":   model,
	})
}

func (c *CrudController) Update(ctx *gin.Context) {
	model := c.getModel()
	ctx.JSON(200, gin.H{
		"message": "Update",
		"model":   model,
	})
}

func (c *CrudController) Delete(ctx *gin.Context) {
	model := c.getModel()
	ctx.JSON(200, gin.H{
		"message": "Delete",
		"model":   model,
	})
}

package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserController struct{}

func (u UserController) Retrieve(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Retrieve user",
	})
}

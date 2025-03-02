package server

import (
	"github.com/BlenDMinh/dutgrad-server/controller"
	"github.com/gin-gonic/gin"
)

func GetRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	v1 := router.Group("/v1")
	{
		v1.GET("", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "Hello, World!",
			})
		})

		userController := new(controller.UserController)
		userGroup := v1.Group("/user")
		{
			userController.Register(userGroup)
		}
	}

	return router
}

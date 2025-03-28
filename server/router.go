package server

import (
	"github.com/BlenDMinh/dutgrad-server/controllers"
	"github.com/BlenDMinh/dutgrad-server/middlewares"
	"github.com/gin-gonic/gin"
)

func GetRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.HandleMethodNotAllowed = true

	v1 := router.Group("/v1")
	{
		v1.GET("", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "Hello, World!",
			})
		})

		userController := controllers.NewUserController()
		userGroup := v1.Group("/user")
		{
			userGroup.GET("/me", middlewares.AuthMiddleware(), userController.GetCurrentUser)
			userController.RegisterCRUD(userGroup)
		}

		authGroup := v1.Group("/auth")
		{
			authController := controllers.NewAuthController()
			authGroup.POST("/register", authController.Register)
			authGroup.POST("/login", authController.Login)
			authGroup.POST("/external-auth", authController.ExternalAuth)
			authGroup.POST("/exchange-state", authController.ExchangeState)

			oauthGroup := authGroup.Group("/oauth")
			{
				oauthController := controllers.NewOAuthController()
				oauthGroup.GET("/google", oauthController.GoogleOAuth)
			}
		}

		documentController := controllers.NewDocumentController()
		documentGroup := v1.Group("/documents")
		{
			documentController.RegisterCRUD(documentGroup)
		}

		spaceController := controllers.NewSpaceController()
		spaceGroup := v1.Group("/spaces")
		{
			spaceController.RegisterCRUD(spaceGroup)
		}

		spaceDocumentsGroup := v1.Group("space/:space_id/documents")
		{
			spaceDocumentsGroup.GET("", documentController.GetBySpaceID)
		}

		userQuerySessionController := controllers.NewUserQuerySessionController()
		userQuerySessionGroup := v1.Group("/user-query-sessions")
		{
			userQuerySessionController.RegisterCRUD(userQuerySessionGroup)
		}

	}

	return router
}

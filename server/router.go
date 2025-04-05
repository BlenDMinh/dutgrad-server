package server

import (
	"time"

	"github.com/BlenDMinh/dutgrad-server/configs"
	"github.com/BlenDMinh/dutgrad-server/controllers"
	"github.com/BlenDMinh/dutgrad-server/middlewares"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func GetRouter() *gin.Engine {
	env := configs.GetEnv()
	router := gin.New()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     env.AllowOrigins,
		AllowMethods:     []string{"GET", "PUT", "PATCH", "POST", "DELETE"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
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
			spaceGroup.GET("", spaceController.Retrieve)
			spaceGroup.POST("", middlewares.AuthMiddleware(), spaceController.CreateSpace)
			spaceGroup.GET("/:id", spaceController.RetrieveOne)
			spaceGroup.PUT("/:id", spaceController.Update)
			spaceGroup.PATCH("/:id", spaceController.Patch)
			spaceGroup.DELETE("/:id", spaceController.Delete)

			spaceGroup.GET("/public", spaceController.GetPublicSpaces)
			spaceGroup.GET("/me", middlewares.AuthMiddleware(), userController.GetMySpaces)
			spaceGroup.GET("/user/:user_id", userController.GetUserSpaces)
		}

		spaceInvitationController := controllers.NewSpaceInvitationController()
		spaceInvitationGroup := v1.Group("/space-invitations")
		{
			spaceInvitationController.RegisterCRUD(spaceInvitationGroup)
		}

		spaceDocumentsGroup := v1.Group("space/:space_id/documents")
		{
			spaceDocumentsGroup.GET("", documentController.GetBySpaceID)
		}

		spaceInvitationLinkController := controllers.NewSpaceInvitationLinkController()
		spaceInvitationLinkGroup := v1.Group("/space-invitation-links")
		{
			spaceInvitationLinkController.RegisterCRUD(spaceInvitationLinkGroup)
		}

		userQuerySessionController := controllers.NewUserQuerySessionController()
		userQuerySessionGroup := v1.Group("/user-query-sessions")
		{
			userQuerySessionController.RegisterCRUD(userQuerySessionGroup)
		}

		userQueryController := controllers.NewUserQueryController()
		userQueryGroup := v1.Group("/user-query")
		{
			userQueryController.RegisterCRUD(userQueryGroup)
		}

	}

	return router
}

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

		userInvitationController := controllers.NewUserController()
		userInvitationGroup := v1.Group("")
		{
			userInvitationGroup.GET("invitations/me", middlewares.AuthMiddleware(), userInvitationController.GetMyInvitations)
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
			documentGroup.GET("", documentController.Retrieve)
			documentGroup.GET("/:id", documentController.RetrieveOne)
			documentGroup.PUT("/:id", documentController.Update)
			documentGroup.PATCH("/:id", documentController.Patch)
			// documentGroup.DELETE("/:id", documentController.Delete)
			documentGroup.POST("/upload", middlewares.AuthMiddleware(), documentController.UploadDocument)
			documentGroup.DELETE("/:id", middlewares.AuthMiddleware(), documentController.DeleteDocument)
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
			spaceGroup.GET("/:id/members", spaceController.GetMembers)
			spaceGroup.GET("/:id/invitations", spaceController.GetInvitations)
			spaceGroup.PUT("/:id/invitation-link", middlewares.AuthMiddleware(), spaceController.GetInvitationLink)
			spaceGroup.POST("/:id/invitations", middlewares.AuthMiddleware(), spaceController.InviteUserToSpace)
			spaceGroup.GET("/roles", spaceController.GetSpaceRoles)
			spaceGroup.GET("/:id/user-role", middlewares.AuthMiddleware(), spaceController.GetUserRole)
			spaceGroup.POST("/join", middlewares.AuthMiddleware(), spaceController.JoinSpace)
			spaceGroup.POST("/:id/join-public", middlewares.AuthMiddleware(), spaceController.JoinPublicSpace)
			spaceGroup.GET("/public", spaceController.GetPublicSpaces)
			spaceGroup.GET("/me", middlewares.AuthMiddleware(), userController.GetMySpaces)
			spaceGroup.GET("/user/:user_id", userController.GetUserSpaces)
		}

		spaceInvitationController := controllers.NewSpaceInvitationController()
		spaceInvitationGroup := v1.Group("/space-invitations")
		{
			spaceInvitationGroup.GET("", spaceInvitationController.Retrieve)
			spaceInvitationGroup.GET("/:id", spaceInvitationController.RetrieveOne)
			spaceInvitationGroup.PUT("/:id", spaceInvitationController.Update)
			spaceInvitationGroup.PATCH("/:id", spaceInvitationController.Patch)
			spaceInvitationGroup.DELETE("/:id", spaceInvitationController.Delete)
			spaceInvitationGroup.PUT("/:id/accept", middlewares.AuthMiddleware(), spaceInvitationController.AcceptInvitation)
			spaceInvitationGroup.PUT("/:id/reject", middlewares.AuthMiddleware(), spaceInvitationController.RejectInvitation)
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
			userQuerySessionGroup.POST("/begin-chat-session", middlewares.AuthMiddleware(), userQuerySessionController.BeginChatSession)
		}

		userQueryController := controllers.NewUserQueryController()
		userQueryGroup := v1.Group("/user-query")
		{
			userQueryController.RegisterCRUD(userQueryGroup)
			userQueryGroup.POST("/ask", middlewares.AuthMiddleware(), userQueryController.Ask)
		}
	}

	return router
}

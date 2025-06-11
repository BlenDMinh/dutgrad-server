package server

import (
	"time"

	"github.com/BlenDMinh/dutgrad-server/configs"
	"github.com/BlenDMinh/dutgrad-server/controllers"
	"github.com/BlenDMinh/dutgrad-server/middlewares"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func GetRouter(
	userController *controllers.UserController,
	authController *controllers.AuthController,
	oauthController *controllers.OAuthController,
	documentController *controllers.DocumentController,
	spaceController *controllers.SpaceController,
	spaceInvitationController *controllers.SpaceInvitationController,
	spaceInvitationLinkController *controllers.SpaceInvitationLinkController,
	userQuerySessionController *controllers.UserQuerySessionController,
	userQueryController *controllers.UserQueryController,
	spaceApiKeyController *controllers.SpaceApiKeyController,
	chatRateLimiter gin.HandlerFunc,
) *gin.Engine {
	env := configs.GetEnv()
	router := gin.New()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     env.AllowOrigins,
		AllowMethods:     []string{"GET", "PUT", "PATCH", "POST", "DELETE", "OPTIONS", "HEAD"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
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

		userGroup := v1.Group("/user")
		{
			userGroup.GET("/me", middlewares.AuthMiddleware(), userController.GetCurrentUser)
			userGroup.GET("/search", middlewares.AuthMiddleware(), userController.SearchUsers)
			userGroup.GET("/tier", middlewares.AuthMiddleware(), userController.GetUserTier)
			userGroup.GET("/auth-method", middlewares.AuthMiddleware(), userController.GetUserAuthMethod)
			userGroup.PATCH("/password", middlewares.AuthMiddleware(), userController.UpdatePassword)
			userController.RegisterCRUD(userGroup)
		}

		userInvitationGroup := v1.Group("")
		{
			userInvitationGroup.GET("invitations/me", middlewares.AuthMiddleware(), userController.GetMyInvitations)
		}

		authGroup := v1.Group("/auth")
		{
			authGroup.POST("/register", authController.Register)
			authGroup.POST("/login", authController.Login)
			authGroup.POST("/external-auth", authController.ExternalAuth)
			authGroup.POST("/exchange-state", authController.ExchangeState)
			authGroup.POST("/verify-mfa", authController.VerifyMFA)

			mfaManagementGroup := authGroup.Group("/mfa")
			mfaManagementGroup.Use(middlewares.AuthMiddleware())
			{
				mfaManagementGroup.GET("/status", authController.GetMFAStatus)

				mfaManagementGroup.POST("/setup", authController.SetupMFA)
				mfaManagementGroup.POST("/verify", authController.ConfirmMFA)
				mfaManagementGroup.POST("/disable", authController.DisableMFA)
			}

			oauthGroup := authGroup.Group("/oauth")
			{
				oauthGroup.GET("/google", oauthController.GoogleOAuth)
			}
		}

		documentGroup := v1.Group("/documents")
		{
			documentGroup.GET("", documentController.Retrieve)
			documentGroup.GET("/:id", documentController.RetrieveOne)

			documentGroup.HEAD("/count/me", middlewares.AuthMiddleware(), documentController.GetUserDocumentCount)

			documentGroup.POST("/upload", middlewares.AuthMiddleware(), documentController.UploadDocument)

			documentGroup.PUT("/:id", documentController.Update)

			documentGroup.PATCH("/:id", documentController.Patch)

			documentGroup.DELETE("/:id", middlewares.AuthMiddleware(), documentController.DeleteDocument)
		}

		spaceGroup := v1.Group("/spaces")
		{
			spaceGroup.GET("", spaceController.Retrieve)
			spaceGroup.GET("/roles", spaceController.GetSpaceRoles)
			spaceGroup.GET("/public", spaceController.GetPublicSpaces)
			spaceGroup.GET("/popular", spaceController.GetPopularSpaces)
			spaceGroup.GET("/user/:id", userController.GetUserSpaces)
			spaceGroup.GET("/me", middlewares.AuthMiddleware(), userController.GetMySpaces)

			spaceGroup.HEAD("/count/me", middlewares.AuthMiddleware(), spaceController.CountMySpaces)

			spaceGroup.POST("/join", middlewares.AuthMiddleware(), spaceController.JoinSpace)
			spaceGroup.POST("", middlewares.AuthMiddleware(), spaceController.CreateSpace)

			detailGroup := spaceGroup.Group("/:id")
			detailGroup.Use(middlewares.AuthMiddleware())
			{
				detailGroup.GET("", spaceController.RetrieveOne)
				detailGroup.PUT("", spaceController.Update)
				detailGroup.PATCH("", spaceController.Patch)
				detailGroup.DELETE("", spaceController.Delete)

				detailGroup.GET("/members", spaceController.GetMembers)
				detailGroup.GET("/members/count", spaceController.CountSpaceMembers)
				detailGroup.GET("/invitations", spaceController.GetInvitations)
				detailGroup.GET("/user-role", spaceController.GetUserRole)
				detailGroup.GET("/documents", documentController.GetBySpaceID)

				detailGroup.PUT("/invitation-link", spaceController.GetInvitationLink)

				detailGroup.POST("/invitations", spaceController.InviteUserToSpace)
				detailGroup.POST("/join-public", spaceController.JoinPublicSpace)

				detailGroup.PATCH("/members/:memberId/role", middlewares.AuthMiddleware(), spaceController.UpdateUserRole)

				detailGroup.DELETE("/members/:memberId", middlewares.AuthMiddleware(), spaceController.RemoveMember)
				apiKeyGroup := detailGroup.Group("/api-keys")
				{
					apiKeyGroup.GET("", spaceApiKeyController.List)
					apiKeyGroup.GET("/:keyId", spaceApiKeyController.GetOne)

					apiKeyGroup.POST("", spaceApiKeyController.Create)

					apiKeyGroup.DELETE("/:keyId", spaceApiKeyController.Delete)
				}
			}
			spaceGroup.POST("/:id/chat", middlewares.RequireApiKey(), spaceController.Chat)
		}
		spaceInvitationGroup := v1.Group("/space-invitations")
		{
			spaceInvitationGroup.GET("/count", middlewares.AuthMiddleware(), spaceInvitationController.GetInvitationCount)
			spaceInvitationGroup.GET("", spaceInvitationController.Retrieve)
			spaceInvitationGroup.GET("/:id", spaceInvitationController.RetrieveOne)

			spaceInvitationGroup.PUT("/:id", spaceInvitationController.Update)
			spaceInvitationGroup.PUT("/:id/accept", middlewares.AuthMiddleware(), spaceInvitationController.AcceptInvitation)
			spaceInvitationGroup.PUT("/:id/reject", middlewares.AuthMiddleware(), spaceInvitationController.RejectInvitation)

			spaceInvitationGroup.PATCH("/:id", spaceInvitationController.Patch)

			spaceInvitationGroup.DELETE("/:id", spaceInvitationController.Delete)
		}

		spaceInvitationLinkGroup := v1.Group("/space-invitation-links")
		{
			spaceInvitationLinkController.RegisterCRUD(spaceInvitationLinkGroup)
		}

		userQuerySessionGroup := v1.Group("/user-query-sessions")
		userQuerySessionGroup.Use(middlewares.AuthMiddleware())
		{
			userQuerySessionController.RegisterCRUD(userQuerySessionGroup)
			userQuerySessionGroup.GET("/me", userQuerySessionController.GetMyChatSessions)
			userQuerySessionGroup.GET("/:id/temp-message", userQuerySessionController.GetTempMessageByID)
			userQuerySessionGroup.GET("/:id/history", userQuerySessionController.GetChatHistory)

			userQuerySessionGroup.HEAD("/me", userQuerySessionController.CountMyChatSessions)

			userQuerySessionGroup.POST("/begin-chat-session", userQuerySessionController.BeginChatSession)

			userQuerySessionGroup.DELETE("/:id/history", userQuerySessionController.ClearChatHistory)
		}

		userQueryGroup := v1.Group("/user-query")
		userQueryGroup.Use(middlewares.AuthMiddleware())
		{
			userQueryController.RegisterCRUD(userQueryGroup)

			userQueryGroup.POST("/ask", chatRateLimiter, userQueryController.Ask)
		}
	}

	return router
}

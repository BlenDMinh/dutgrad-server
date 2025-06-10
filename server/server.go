package server

import (
	"strconv"

	"github.com/BlenDMinh/dutgrad-server/configs"
	"github.com/BlenDMinh/dutgrad-server/controllers"
	"github.com/BlenDMinh/dutgrad-server/databases/repositories"
	"github.com/BlenDMinh/dutgrad-server/middlewares"
	"github.com/BlenDMinh/dutgrad-server/services"
	"github.com/BlenDMinh/dutgrad-server/services/oauth"
	"github.com/BlenDMinh/dutgrad-server/services/oauth/providers"
)

func Init() {
	// Repository initialization
	userRepo := repositories.NewUserRepository()
	userMFARepo := repositories.NewUserMFARepository()
	authCredentialRepo := repositories.NewUserAuthCredentialRepository()
	spaceInvitationRepo := repositories.NewSpaceInvitationRepository()
	spaceInvitationLinkRepo := repositories.NewSpaceInvitationLinkRepository()
	documentRepo := repositories.NewDocumentRepository()

	// External service initialization
	ragServerService := services.NewRAGServerService()
	// redisService := services.NewRedisService()
	memoryStorage := services.NewInMemoryStorage()
	mfaService := services.NewMFAService(
		memoryStorage,
		userRepo,
		userMFARepo,
		authCredentialRepo,
	)

	// Service initialization
	userService := services.NewUserService()
	authService := services.NewAuthService()
	documentService := services.NewDocumentService(ragServerService)
	spaceService := services.NewSpaceService(
		spaceInvitationLinkRepo,
		ragServerService,
		userRepo,
		spaceInvitationRepo,
		documentRepo,
	)
	spaceInvitationService := services.NewSpaceInvitationService()
	spaceInvitationLinkService := services.NewSpaceInvitationLinkService()
	userQuerySessionService := services.NewUserQuerySessionService()
	userQueryService := services.NewUserQueryService()
	spaceApiKeyService := services.NewSpaceApiKeyService()

	// Controller initialization
	userController := controllers.NewUserController(userService)
	authController := controllers.NewAuthController(
		authService,
		userService,
		memoryStorage,
		mfaService,
	)
	providers := map[string]oauth.OAuthProvider{
		"google": providers.NewGoogleOAuthProvider(),
	}
	oauthController := controllers.NewOAuthController(
		providers,
		authService,
		memoryStorage,
		mfaService,
	)
	documentController := controllers.NewDocumentController(documentService, spaceService)
	spaceController := controllers.NewSpaceController(spaceService)
	spaceInvitationController := controllers.NewSpaceInvitationController(spaceInvitationService)
	spaceInvitationLinkController := controllers.NewSpaceInvitationLinkController(spaceInvitationLinkService)
	userQuerySessionController := controllers.NewUserQuerySessionController(userQuerySessionService)
	userQueryController := controllers.NewUserQueryController(userQueryService)
	spaceApiKeyController := controllers.NewSpaceApiKeyController(spaceApiKeyService)

	config := configs.GetEnv()

	// Middleware initialization
	chatRateLimiter := middlewares.ChatRateLimiter(userService)

	// Router initialization
	r := GetRouter(
		userController,
		authController,
		oauthController,
		documentController,
		spaceController,
		spaceInvitationController,
		spaceInvitationLinkController,
		userQuerySessionController,
		userQueryController,
		spaceApiKeyController,
		chatRateLimiter,
	)

	r.Run(":" + strconv.Itoa(config.Port))
}

func Close() {

}

package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/BlenDMinh/dutgrad-server/configs"
	"github.com/BlenDMinh/dutgrad-server/models/dtos"
	"github.com/BlenDMinh/dutgrad-server/services"
	"github.com/BlenDMinh/dutgrad-server/services/oauth"
	"github.com/BlenDMinh/dutgrad-server/services/oauth/providers"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	StateTokenExpiration = 10 * time.Minute
	MFATokenExpiration   = 5 * time.Minute
)

type OAuthController struct {
	providers    map[string]oauth.OAuthProvider
	authService  *services.AuthService
	redisService *services.RedisService
	mfaService   *services.MFAService
}

func NewOAuthController() *OAuthController {
	providers := map[string]oauth.OAuthProvider{
		"google": providers.NewGoogleOAuthProvider(),
	}

	return &OAuthController{
		providers:    providers,
		authService:  &services.AuthService{},
		redisService: services.NewRedisService(),
		mfaService:   services.NewMFAService(),
	}
}

func (c *OAuthController) HandleOAuthCallback(ctx *gin.Context, providerName string) {
	provider, exists := c.providers[providerName]
	if !exists {
		ctx.Redirect(http.StatusTemporaryRedirect,
			fmt.Sprintf("%s/auth/error?code=invalid_provider&message=%s",
				configs.GetEnv().WebClientURL,
				"Provider not supported"))
		return
	}

	code := ctx.Query("code")
	if code == "" {
		ctx.Redirect(http.StatusTemporaryRedirect,
			fmt.Sprintf("%s/auth/error?code=missing_code&message=%s",
				configs.GetEnv().WebClientURL,
				"Authorization code missing"))
		return
	}

	token, err := provider.GetConfig().Exchange(context.Background(), code)
	if err != nil {
		log.Printf("Code exchange error: %v", err)
		ctx.Redirect(http.StatusTemporaryRedirect,
			fmt.Sprintf("%s/auth/error?code=exchange_failed&message=%s",
				configs.GetEnv().WebClientURL,
				"Failed to exchange authorization code"))
		return
	}

	userInfo, err := provider.GetUserInfo(token)
	if err != nil {
		log.Printf("Failed to get user info: %v", err)
		ctx.Redirect(http.StatusTemporaryRedirect,
			fmt.Sprintf("%s/auth/error?code=user_info_failed&message=%s",
				configs.GetEnv().WebClientURL,
				"Failed to retrieve user information"))
		return
	}

	externalAuthDto := dtos.ExternalAuthDTO{
		Email:      userInfo.Email,
		Username:   userInfo.Username,
		ExternalID: userInfo.ID,
		AuthType:   userInfo.Provider,
	}

	user, jwt_token, expiresAt, IsNewUser, err := c.authService.ExternalAuth(&externalAuthDto)
	if err != nil {
		log.Printf("Authentication error: %v", err)
		ctx.Redirect(http.StatusTemporaryRedirect,
			fmt.Sprintf("%s/auth/error?code=auth_failed&message=%s",
				configs.GetEnv().WebClientURL,
				"Authentication failed"))
		return
	}

	mfaEnabled, err := c.mfaService.GetUserMFAStatus(user.ID)
	if err != nil {
		log.Printf("Failed to check MFA status: %v", err)
		ctx.Redirect(http.StatusTemporaryRedirect,
			fmt.Sprintf("%s/auth/error?code=mfa_check_failed&message=%s",
				configs.GetEnv().WebClientURL,
				"Failed to check MFA status"))
		return
	}

	if mfaEnabled {
		tempToken, _, err := c.mfaService.CreateTempToken(user.ID)
		if err != nil {
			log.Printf("Failed to create MFA temp token: %v", err)
			ctx.Redirect(http.StatusTemporaryRedirect,
				fmt.Sprintf("%s/auth/error?code=mfa_token_failed&message=%s",
					configs.GetEnv().WebClientURL,
					"Failed to create MFA token"))
			return
		}

		if err := c.redisService.Set(tempToken, user.ID, MFATokenExpiration); err != nil {
			log.Printf("Redis error: %v", err)
			ctx.Redirect(http.StatusTemporaryRedirect,
				fmt.Sprintf("%s/auth/error?code=redis_error&message=%s",
					configs.GetEnv().WebClientURL,
					"Internal server error"))
			return
		}

		mfaURL := fmt.Sprintf("%s/auth/mfa?state=%s",
			configs.GetEnv().WebClientURL,
			tempToken)

		ctx.Redirect(http.StatusTemporaryRedirect, mfaURL)
		return
	}

	stateToken := uuid.New().String()
	authResponse := dtos.AuthResponse{
		Token:     jwt_token,
		User:      user,
		IsNewUser: IsNewUser,
		Expires:   expiresAt,
	}

	if err := c.redisService.Set(stateToken, authResponse, StateTokenExpiration); err != nil {
		log.Printf("Redis error: %v", err)
		ctx.Redirect(http.StatusTemporaryRedirect,
			fmt.Sprintf("%s/auth/error?code=redis_error&message=%s",
				configs.GetEnv().WebClientURL,
				"Internal server error"))
		return
	}

	successURL := fmt.Sprintf("%s/auth/success?state=%s",
		configs.GetEnv().WebClientURL,
		stateToken)

	ctx.Redirect(http.StatusTemporaryRedirect, successURL)
}

func (c *OAuthController) GoogleOAuth(ctx *gin.Context) {
	c.HandleOAuthCallback(ctx, "google")
}

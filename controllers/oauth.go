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

type OAuthController struct {
	providers    map[string]oauth.OAuthProvider
	authService  *services.AuthService
	redisService *services.RedisService
}

func NewOAuthController() *OAuthController {
	providers := map[string]oauth.OAuthProvider{
		"google": providers.NewGoogleOAuthProvider(),
	}

	return &OAuthController{
		providers:    providers,
		authService:  &services.AuthService{},
		redisService: services.NewRedisService(),
	}
}

func (c *OAuthController) HandleOAuthCallback(ctx *gin.Context, providerName string) {
	provider, exists := c.providers[providerName]
	if !exists {
		ctx.Redirect(http.StatusTemporaryRedirect,
			fmt.Sprintf("%s/auth/error", configs.GetEnv().WebClientURL))
		return
	}

	code := ctx.Query("code")
	if code == "" {
		ctx.Redirect(http.StatusTemporaryRedirect,
			fmt.Sprintf("%s/auth/error", configs.GetEnv().WebClientURL))
		return
	}

	token, err := provider.GetConfig().Exchange(context.Background(), code)
	if err != nil {
		log.Printf("Code exchange error: %v", err)
		ctx.Redirect(http.StatusTemporaryRedirect,
			fmt.Sprintf("%s/auth/error", configs.GetEnv().WebClientURL))
		return
	}

	userInfo, err := provider.GetUserInfo(token)
	if err != nil {
		log.Printf("Failed to get user info: %v", err)
		ctx.Redirect(http.StatusTemporaryRedirect,
			fmt.Sprintf("%s/auth/error", configs.GetEnv().WebClientURL))
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
			fmt.Sprintf("%s/auth/error", configs.GetEnv().WebClientURL))
		return
	}

	stateToken := uuid.New().String()
	authResponse := dtos.AuthResponse{
		Token:     jwt_token,
		User:      user,
		IsNewUser: IsNewUser,
		Expires:   expiresAt,
	}

	if err := c.redisService.Set(stateToken, authResponse, 5*time.Minute); err != nil {
		log.Printf("Redis error: %v", err)
		ctx.Redirect(http.StatusTemporaryRedirect,
			fmt.Sprintf("%s/auth/error", configs.GetEnv().WebClientURL))
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

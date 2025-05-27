package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/BlenDMinh/dutgrad-server/configs"
	"github.com/BlenDMinh/dutgrad-server/models/dtos"
	"github.com/BlenDMinh/dutgrad-server/services"
	"github.com/BlenDMinh/dutgrad-server/services/oauth"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	StateTokenExpiration = 10 * time.Minute
	MFATokenExpiration   = 5 * time.Minute
)

func redirectToError(ctx *gin.Context, code string, message string) {
	errorURL := fmt.Sprintf("%s/auth/error?code=%s&message=%s",
		configs.GetEnv().WebClientURL,
		url.QueryEscape(code),
		url.QueryEscape(message))
	ctx.Redirect(http.StatusTemporaryRedirect, errorURL)
}

type OAuthController struct {
	providers   map[string]oauth.OAuthProvider
	authService *services.AuthService
	kvStorage   services.KVStorage
	mfaService  *services.MFAService
}

func NewOAuthController(
	providers map[string]oauth.OAuthProvider,
	authService *services.AuthService,
	kvStorage services.KVStorage,
	mfaService *services.MFAService,
) *OAuthController {
	return &OAuthController{
		providers:   providers,
		authService: authService,
		kvStorage:   kvStorage,
		mfaService:  mfaService,
	}
}

func (c *OAuthController) HandleOAuthCallback(ctx *gin.Context, providerName string) {
	provider, exists := c.providers[providerName]
	if !exists {
		redirectToError(ctx, "invalid_provider", "Provider not supported")
		return
	}

	code := ctx.Query("code")
	if code == "" {
		redirectToError(ctx, "missing_code", "Authorization code missing")
		return
	}

	token, err := provider.GetConfig().Exchange(context.Background(), code)
	if err != nil {
		log.Printf("Code exchange error: %v", err)
		redirectToError(ctx, "exchange_failed", "Failed to exchange authorization code")
		return
	}

	userInfo, err := provider.GetUserInfo(token)
	if err != nil {
		log.Printf("Failed to get user info: %v", err)
		redirectToError(ctx, "user_info_failed", "Failed to retrieve user information")
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
		redirectToError(ctx, "auth_failed", "Authentication failed")
		return
	}

	mfaEnabled, err := c.mfaService.GetUserMFAStatus(user.ID)
	if err != nil {
		log.Printf("Failed to check MFA status: %v", err)
		redirectToError(ctx, "mfa_check_failed", "Failed to check MFA status")
		return
	}

	if mfaEnabled {
		tempToken, _, err := c.mfaService.CreateTempToken(user.ID)
		if err != nil {
			log.Printf("Failed to create MFA temp token: %v", err)
			redirectToError(ctx, "mfa_token_failed", "Failed to create MFA token")
			return
		}

		if err := c.kvStorage.Set(tempToken, user.ID, MFATokenExpiration); err != nil {
			log.Printf("Redis error: %v", err)
			redirectToError(ctx, "redis_error", "Internal server error")
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

	if err := c.kvStorage.Set(stateToken, authResponse, StateTokenExpiration); err != nil {
		log.Printf("Redis error: %v", err)
		redirectToError(ctx, "redis_error", "Internal server error")
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

func (c *OAuthController) ExchangeState(ctx *gin.Context) {
	state := ctx.Query("state")
	if state == "" {
		HandleError(ctx, http.StatusBadRequest, "Invalid state token", nil)
		return
	}

	authDataJSON, err := c.kvStorage.Get(state)
	if err != nil {
		HandleError(ctx, http.StatusNotFound, "State token expired or invalid", nil)
		return
	}

	c.kvStorage.Delete(state)

	var authResponse dtos.AuthResponse
	if err := json.Unmarshal([]byte(authDataJSON), &authResponse); err != nil {
		HandleError(ctx, http.StatusInternalServerError, "Failed to parse auth data", err)
		return
	}

	HandleSuccess(ctx, "Token exchange successful", authResponse)
}

func (c *OAuthController) VerifyOAuthMFA(ctx *gin.Context) {
	var req dtos.MFALoginCompleteRequest
	if !HandleBindJSON(ctx, &req) {
		return
	}

	tempToken := ctx.Query("state")
	if tempToken == "" {
		HandleError(ctx, http.StatusBadRequest, "Missing temporary token", nil)
		return
	}

	userIDStr, err := c.kvStorage.Get(tempToken)
	if err != nil {
		HandleError(ctx, http.StatusUnauthorized, "Invalid or expired temporary token", err)
		return
	}

	var userID uint
	if err := json.Unmarshal([]byte(userIDStr), &userID); err != nil {
		HandleError(ctx, http.StatusInternalServerError, "Failed to parse user ID", err)
		return
	}

	isValid := c.mfaService.VerifyMFACode(userID, req.Code, req.UseBackupCode)
	if !isValid {
		HandleError(ctx, http.StatusUnauthorized, "Invalid MFA code", nil)
		return
	}

	user, err := services.NewUserService().GetById(userID)
	if err != nil {
		HandleError(ctx, http.StatusInternalServerError, "Failed to get user", err)
		return
	}

	token, expiresAt, err := c.mfaService.CompleteLogin(userID)
	if err != nil {
		HandleError(ctx, http.StatusInternalServerError, "Failed to complete login", err)
		return
	}

	c.kvStorage.Delete(tempToken)

	HandleSuccess(ctx, "Login successful", dtos.AuthResponse{
		Token:   token,
		User:    user,
		Expires: expiresAt,
	})
}

func (c *OAuthController) GetAuthorizationURL(ctx *gin.Context, providerName string) {
	provider, exists := c.providers[providerName]
	if !exists {
		HandleError(ctx, http.StatusBadRequest, "Provider not supported", nil)
		return
	}

	state := uuid.New().String()
	authURL := provider.GetConfig().AuthCodeURL(state)

	HandleSuccess(ctx, "Authorization URL generated", gin.H{
		"url":   authURL,
		"state": state,
	})
}

func (c *OAuthController) GoogleAuthURL(ctx *gin.Context) {
	c.GetAuthorizationURL(ctx, "google")
}

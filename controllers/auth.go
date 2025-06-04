package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/BlenDMinh/dutgrad-server/models/dtos"
	"github.com/BlenDMinh/dutgrad-server/services"
	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authService *services.AuthService
	userService services.UserService
	kvStorage   services.KVStorage
	mfaService  *services.MFAService
}

func NewAuthController(
	authService *services.AuthService,
	userService services.UserService,
	kvStorage services.KVStorage,
	mfaService *services.MFAService,
) *AuthController {
	return &AuthController{
		authService: authService,
		userService: userService,
		kvStorage:   kvStorage,
		mfaService:  mfaService,
	}
}

func (ac *AuthController) Register(ctx *gin.Context) {
	var req dtos.RegisterRequest
	if !HandleBindJSON(ctx, &req) {
		return
	}

	dto := dtos.RegisterDTO(req)
	user, token, expiresAt, err := ac.authService.RegisterUser(&dto)

	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "user with this email already exists" {
			statusCode = http.StatusConflict
		}
		HandleError(ctx, statusCode, "Registration failed", err)
		return
	}

	HandleCreated(ctx, "User registered successfully", dtos.AuthResponse{
		Token:     token,
		User:      user,
		IsNewUser: true,
		Expires:   expiresAt,
	})
}

func (ac *AuthController) Login(ctx *gin.Context) {
	var req dtos.LoginRequest
	if !HandleBindJSON(ctx, &req) {
		return
	}

	user, requiresMFA, err := ac.mfaService.FirstFactorAuth(req.Email, req.Password)
	if err != nil {
		HandleError(ctx, http.StatusUnauthorized, "Authentication failed", err)
		return
	}

	if !requiresMFA {
		token, expiresAt, err := ac.mfaService.CompleteLogin(user.ID)
		if err != nil {
			HandleError(ctx, http.StatusInternalServerError, "Login failed", err)
			return
		}

		HandleSuccess(ctx, "Login successful", dtos.AuthResponse{
			Token:   token,
			User:    user,
			Expires: expiresAt,
		})
		return
	}

	tempToken, expiresAt, err := ac.mfaService.CreateTempToken(user.ID)
	if err != nil {
		HandleError(ctx, http.StatusInternalServerError, "Failed to create temporary token", err)
		return
	}

	HandleSuccess(ctx, "MFA verification required", gin.H{
		"requires_mfa": true,
		"temp_token":   tempToken,
		"expires_at":   expiresAt.Format(time.RFC3339),
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		},
	})
}

func (ac *AuthController) VerifyMFA(ctx *gin.Context) {
	var req dtos.MFALoginCompleteRequest
	if !HandleBindJSON(ctx, &req) {
		return
	}

	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		HandleError(ctx, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	tempToken := strings.TrimPrefix(authHeader, "Bearer ")
	log.Printf("Temp token: %s", tempToken)

	userID, err := ac.mfaService.GetUserIDFromTempToken(tempToken)
	log.Printf("User ID from temp token: %d", userID)
	if err != nil {
		HandleError(ctx, http.StatusUnauthorized, "Invalid or expired temporary token", err)
		return
	}

	isValid := ac.mfaService.VerifyMFACode(userID, req.Code, req.UseBackupCode)
	log.Printf("MFA code valid: %v", isValid)
	if !isValid {
		HandleError(ctx, http.StatusUnauthorized, "Invalid MFA code", nil)
		return
	}

	user, err := ac.userService.GetById(userID)
	if err != nil {
		HandleError(ctx, http.StatusInternalServerError, "Failed to get user", err)
		return
	}

	token, expiresAt, err := ac.mfaService.CompleteLogin(userID)
	if err != nil {
		HandleError(ctx, http.StatusInternalServerError, "Failed to complete login", err)
		return
	}

	HandleSuccess(ctx, "Login successful", dtos.AuthResponse{
		Token:   token,
		User:    user,
		Expires: expiresAt,
	})
}

func (ac *AuthController) GetMFAStatus(ctx *gin.Context) {
	userID, ok := ExtractID(ctx, "user_id")
	if !ok {
		return
	}

	mfaEnabled, err := ac.mfaService.GetUserMFAStatus(userID)
	if err != nil {
		HandleError(ctx, http.StatusInternalServerError, "Failed to get MFA status", err)
		return
	}

	HandleSuccess(ctx, "MFA status retrieved", gin.H{
		"mfa_enabled": mfaEnabled,
	})
}

func (ac *AuthController) SetupMFA(ctx *gin.Context) {
	userID, ok := ExtractID(ctx, "user_id")
	if !ok {
		return
	}

	setupResponse, err := ac.mfaService.GenerateMFASetup(userID)
	if err != nil {
		HandleError(ctx, http.StatusBadRequest, "Failed to set up MFA", err)
		return
	}

	HandleSuccess(ctx, "MFA setup initialized", setupResponse)
}

func (ac *AuthController) ConfirmMFA(ctx *gin.Context) {
	userID, ok := ExtractID(ctx, "user_id")
	if !ok {
		return
	}

	var req dtos.MFAVerifyRequest
	if !HandleBindJSON(ctx, &req) {
		return
	}

	if err := ac.mfaService.VerifyMFASetup(userID, req.Code); err != nil {
		HandleError(ctx, http.StatusBadRequest, "Failed to verify MFA setup", err)
		return
	}

	HandleSuccess(ctx, "MFA has been enabled successfully", nil)
}

func (ac *AuthController) DisableMFA(ctx *gin.Context) {
	userID, ok := ExtractID(ctx, "user_id")
	if !ok {
		return
	}

	if err := ac.mfaService.DisableMFA(userID); err != nil {
		HandleError(ctx, http.StatusBadRequest, "Failed to disable MFA", err)
		return
	}

	HandleSuccess(ctx, "MFA has been disabled successfully", nil)
}

func (ac *AuthController) ExternalAuth(ctx *gin.Context) {
	var req dtos.ExternalAuthRequest
	if !HandleBindJSON(ctx, &req) {
		return
	}

	dto := dtos.ExternalAuthDTO{
		Email:      req.Email,
		Username:   req.Username,
		ExternalID: req.ExternalID,
		AuthType:   req.AuthType,
	}

	user, token, expiresAt, isNewUser, err := ac.authService.ExternalAuth(&dto)

	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "invalid authentication type" {
			statusCode = http.StatusBadRequest
		}
		HandleError(ctx, statusCode, "Authentication failed", err)
		return
	}

	if isNewUser {
		HandleCreated(ctx, "User registered successfully", dtos.AuthResponse{
			Token:     token,
			User:      user,
			IsNewUser: isNewUser,
			Expires:   expiresAt,
		})
	} else {
		HandleSuccess(ctx, "Login successful", dtos.AuthResponse{
			Token:     token,
			User:      user,
			IsNewUser: isNewUser,
			Expires:   expiresAt,
		})
	}
}

func (ac *AuthController) ExchangeState(ctx *gin.Context) {
	state := ctx.Query("state")
	if state == "" {
		HandleError(ctx, http.StatusBadRequest, "Invalid state token", nil)
		return
	}

	authDataJSON, err := ac.kvStorage.Get(state)
	if err != nil {
		HandleError(ctx, http.StatusNotFound, "State token expired or invalid", nil)
		return
	}

	ac.kvStorage.Delete(state)

	var authResponse dtos.AuthResponse
	if err := json.Unmarshal([]byte(authDataJSON), &authResponse); err != nil {
		HandleError(ctx, http.StatusInternalServerError, "Failed to parse auth data", err)
		return
	}

	HandleSuccess(ctx, "Token exchange successful", authResponse)
}

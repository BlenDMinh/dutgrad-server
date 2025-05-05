package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/BlenDMinh/dutgrad-server/models"
	"github.com/BlenDMinh/dutgrad-server/models/dtos"
	"github.com/BlenDMinh/dutgrad-server/services"
	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authService  *services.AuthService
	userService  *services.UserService
	redisService *services.RedisService
	mfaService   *services.MFAService
}

func NewAuthController() *AuthController {
	return &AuthController{
		authService:  services.NewAuthService(),
		userService:  services.NewUserService(),
		redisService: services.NewRedisService(),
		mfaService:   services.NewMFAService(),
	}
}

func (ac *AuthController) Register(ctx *gin.Context) {
	var req dtos.RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		errMsg := err.Error()
		ctx.JSON(http.StatusBadRequest, models.NewErrorResponse(http.StatusBadRequest, "Invalid request format", &errMsg))
		return
	}

	dto := dtos.RegisterDTO(req)
	user, token, expiresAt, err := ac.authService.RegisterUser(&dto)

	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "user with this email already exists" {
			statusCode = http.StatusConflict
		}
		errMsg := err.Error()
		ctx.JSON(statusCode, models.NewErrorResponse(statusCode, "Registration failed", &errMsg))
		return
	}

	ctx.JSON(http.StatusCreated, models.NewSuccessResponse(http.StatusCreated, "User registered successfully", dtos.AuthResponse{
		Token:     token,
		User:      user,
		IsNewUser: true,
		Expires:   expiresAt,
	}))
}

func (ac *AuthController) Login(ctx *gin.Context) {
	var req dtos.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		errMsg := err.Error()
		ctx.JSON(http.StatusBadRequest, models.NewErrorResponse(http.StatusBadRequest, "Invalid request format", &errMsg))
		return
	}

	user, requiresMFA, err := ac.mfaService.FirstFactorAuth(req.Email, req.Password)
	if err != nil {
		errMsg := err.Error()
		ctx.JSON(http.StatusUnauthorized, models.NewErrorResponse(http.StatusUnauthorized, "Authentication failed", &errMsg))
		return
	}

	if !requiresMFA {
		token, expiresAt, err := ac.mfaService.CompleteLogin(user.ID)
		if err != nil {
			errMsg := err.Error()
			ctx.JSON(http.StatusInternalServerError, models.NewErrorResponse(http.StatusInternalServerError, "Login failed", &errMsg))
			return
		}

		ctx.JSON(http.StatusOK, models.NewSuccessResponse(http.StatusOK, "Login successful", dtos.AuthResponse{
			Token:   token,
			User:    user,
			Expires: expiresAt,
		}))
		return
	}

	tempToken, expiresAt, err := ac.mfaService.CreateTempToken(user.ID)
	if err != nil {
		errMsg := err.Error()
		ctx.JSON(http.StatusInternalServerError, models.NewErrorResponse(http.StatusInternalServerError, "Failed to create temporary token", &errMsg))
		return
	}

	ctx.JSON(http.StatusOK, models.NewSuccessResponse(http.StatusOK, "MFA verification required", gin.H{
		"requires_mfa": true,
		"temp_token":   tempToken,
		"expires_at":   expiresAt.Format(time.RFC3339),
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		},
	}))
}

func (ac *AuthController) VerifyMFA(ctx *gin.Context) {
	var req dtos.MFALoginCompleteRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		errMsg := err.Error()
		ctx.JSON(http.StatusBadRequest, models.NewErrorResponse(http.StatusBadRequest, "Invalid request format", &errMsg))
		return
	}

	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		ctx.JSON(http.StatusUnauthorized, models.NewErrorResponse(http.StatusUnauthorized, "Unauthorized", nil))
		return
	}

	tempToken := strings.TrimPrefix(authHeader, "Bearer ")
	log.Printf("Temp token: %s", tempToken)

	userID, err := ac.mfaService.GetUserIDFromTempToken(tempToken)
	log.Printf("User ID from temp token: %d", userID)
	if err != nil {
		errMsg := err.Error()
		ctx.JSON(http.StatusUnauthorized, models.NewErrorResponse(http.StatusUnauthorized, "Invalid or expired temporary token", &errMsg))
		return
	}

	isValid := ac.mfaService.VerifyMFACode(userID, req.Code, req.UseBackupCode)
	log.Printf("MFA code valid: %v", isValid)
	if !isValid {
		ctx.JSON(http.StatusUnauthorized, models.NewErrorResponse(http.StatusUnauthorized, "Invalid MFA code", nil))
		return
	}

	user, err := ac.userService.GetById(userID)
	if err != nil {
		errMsg := err.Error()
		ctx.JSON(http.StatusInternalServerError, models.NewErrorResponse(http.StatusInternalServerError, "Failed to get user", &errMsg))
		return
	}

	token, expiresAt, err := ac.mfaService.CompleteLogin(userID)
	if err != nil {
		errMsg := err.Error()
		ctx.JSON(http.StatusInternalServerError, models.NewErrorResponse(http.StatusInternalServerError, "Failed to complete login", &errMsg))
		return
	}

	ctx.JSON(http.StatusOK, models.NewSuccessResponse(http.StatusOK, "Login successful", dtos.AuthResponse{
		Token:   token,
		User:    user,
		Expires: expiresAt,
	}))
}

func (ac *AuthController) GetMFAStatus(ctx *gin.Context) {
	userID := ctx.GetUint("user_id")

	mfaEnabled, err := ac.mfaService.GetUserMFAStatus(userID)
	if err != nil {
		errMsg := err.Error()
		ctx.JSON(http.StatusInternalServerError, models.NewErrorResponse(http.StatusInternalServerError, "Failed to get MFA status", &errMsg))
		return
	}

	ctx.JSON(http.StatusOK, models.NewSuccessResponse(http.StatusOK, "MFA status retrieved", gin.H{
		"mfa_enabled": mfaEnabled,
	}))
}

func (ac *AuthController) SetupMFA(ctx *gin.Context) {
	userID := ctx.GetUint("user_id")
	setupResponse, err := ac.mfaService.GenerateMFASetup(userID)
	if err != nil {
		errMsg := err.Error()
		ctx.JSON(http.StatusBadRequest, models.NewErrorResponse(http.StatusBadRequest, "Failed to set up MFA", &errMsg))
		return
	}

	ctx.JSON(http.StatusOK, models.NewSuccessResponse(http.StatusOK, "MFA setup initialized", setupResponse))
}

func (ac *AuthController) ConfirmMFA(ctx *gin.Context) {
	userID := ctx.GetUint("user_id")

	var req dtos.MFAVerifyRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		errMsg := err.Error()
		ctx.JSON(http.StatusBadRequest, models.NewErrorResponse(http.StatusBadRequest, "Invalid request format", &errMsg))
		return
	}

	if err := ac.mfaService.VerifyMFASetup(userID, req.Code); err != nil {
		errMsg := err.Error()
		ctx.JSON(http.StatusBadRequest, models.NewErrorResponse(http.StatusBadRequest, "Failed to verify MFA setup", &errMsg))
		return
	}

	ctx.JSON(http.StatusOK, models.NewSuccessResponse(http.StatusOK, "MFA has been enabled successfully", nil))
}

func (ac *AuthController) DisableMFA(ctx *gin.Context) {
	userID := ctx.GetUint("user_id")

	if err := ac.mfaService.DisableMFA(userID); err != nil {
		errMsg := err.Error()
		ctx.JSON(http.StatusBadRequest, models.NewErrorResponse(http.StatusBadRequest, "Failed to disable MFA", &errMsg))
		return
	}

	ctx.JSON(http.StatusOK, models.NewSuccessResponse(http.StatusOK, "MFA has been disabled successfully", nil))
}

func (ac *AuthController) ExternalAuth(ctx *gin.Context) {
	var req dtos.ExternalAuthRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		errMsg := err.Error()
		ctx.JSON(http.StatusBadRequest, models.NewErrorResponse(http.StatusBadRequest, "Invalid request format", &errMsg))
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
		errMsg := err.Error()
		ctx.JSON(statusCode, models.NewErrorResponse(statusCode, "Authentication failed", &errMsg))
		return
	}

	statusCode := http.StatusOK
	message := "Login successful"
	if isNewUser {
		statusCode = http.StatusCreated
		message = "User registered successfully"
	}

	ctx.JSON(statusCode, models.NewSuccessResponse(statusCode, message, dtos.AuthResponse{
		Token:   token,
		User:    user,
		Expires: expiresAt,
	}))
}

func (ac *AuthController) ExchangeState(ctx *gin.Context) {
	state := ctx.Query("state")
	if state == "" {
		ctx.JSON(http.StatusBadRequest, models.NewErrorResponse(
			http.StatusBadRequest,
			"Invalid state token",
			nil))
		return
	}

	authDataJSON, err := ac.redisService.Get(state)
	if err != nil {
		ctx.JSON(http.StatusNotFound, models.NewErrorResponse(
			http.StatusNotFound,
			"State token expired or invalid",
			nil))
		return
	}

	ac.redisService.Del(state)

	var authResponse dtos.AuthResponse
	if err := json.Unmarshal([]byte(authDataJSON), &authResponse); err != nil {
		ctx.JSON(http.StatusInternalServerError, models.NewErrorResponse(
			http.StatusInternalServerError,
			"Failed to parse auth data",
			nil))
		return
	}

	ctx.JSON(http.StatusOK, models.NewSuccessResponse(
		http.StatusOK,
		"Token exchange successful",
		authResponse))
}

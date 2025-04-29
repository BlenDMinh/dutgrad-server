package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/BlenDMinh/dutgrad-server/models"
	"github.com/BlenDMinh/dutgrad-server/models/dtos"
	"github.com/BlenDMinh/dutgrad-server/services"
	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authService  *services.AuthService
	userService  *services.UserService
	redisService *services.RedisService
}

func NewAuthController() *AuthController {
	return &AuthController{
		authService:  &services.AuthService{},
		userService:  &services.UserService{},
		redisService: services.NewRedisService(),
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

	user, token, expiresAt, err := ac.authService.LoginUser(req.Email, req.Password)
	if err != nil {
		errMsg := err.Error()
		ctx.JSON(http.StatusUnauthorized, models.NewErrorResponse(http.StatusUnauthorized, "Authentication failed", &errMsg))
		return
	}

	ctx.JSON(http.StatusOK, models.NewSuccessResponse(http.StatusOK, "Login successful", dtos.AuthResponse{
		Token:   token,
		User:    user,
		Expires: expiresAt,
	}))
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

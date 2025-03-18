package services

import (
	"errors"
	"time"

	"github.com/BlenDMinh/dutgrad-server/databases"
	"github.com/BlenDMinh/dutgrad-server/databases/entities"
	"github.com/BlenDMinh/dutgrad-server/helpers"
	"github.com/BlenDMinh/dutgrad-server/models/dtos"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct{}

// RegisterUser creates a new user with local authentication
func (s *AuthService) RegisterUser(dto *dtos.RegisterDTO) (*entities.User, string, time.Time, error) {
	db := databases.GetDB()

	// Check if user already exists with this email
	var existingUser entities.User
	if err := db.Where("email = ?", dto.Email).First(&existingUser).Error; err == nil {
		return nil, "", time.Time{}, errors.New("user with this email already exists")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, "", time.Time{}, err
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(dto.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", time.Time{}, errors.New("failed to hash password")
	}

	// Create new user and credentials in transaction
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	email := dto.Email
	now := time.Now()
	user := entities.User{
		FirstName: dto.FirstName,
		LastName:  dto.LastName,
		Email:     &email,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		return nil, "", time.Time{}, err
	}

	passwordStr := string(hashedPassword)
	credentials := entities.UserAuthCredential{
		UserID:       user.ID,
		AuthType:     "local",
		PasswordHash: &passwordStr,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := tx.Create(&credentials).Error; err != nil {
		tx.Rollback()
		return nil, "", time.Time{}, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, "", time.Time{}, err
	}

	// Generate JWT token
	token, expiresAt, err := helpers.GenerateJWTToken(user.ID)
	if err != nil {
		return nil, "", time.Time{}, errors.New("failed to generate authentication token")
	}

	return &user, token, expiresAt, nil
}

// LoginUser authenticates a user with email and password
func (s *AuthService) LoginUser(email, password string) (*entities.User, string, time.Time, error) {
	db := databases.GetDB()

	// Find user by email
	var user entities.User
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, "", time.Time{}, errors.New("invalid email or password")
	}

	// Get user credentials
	var credentials entities.UserAuthCredential
	if err := db.Where("user_id = ? AND auth_type = ?", user.ID, "local").First(&credentials).Error; err != nil {
		return nil, "", time.Time{}, errors.New("invalid email or password")
	}

	// Verify password
	if credentials.PasswordHash == nil {
		return nil, "", time.Time{}, errors.New("account has no password set")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(*credentials.PasswordHash), []byte(password)); err != nil {
		return nil, "", time.Time{}, errors.New("invalid email or password")
	}

	// Generate JWT token
	token, expiresAt, err := helpers.GenerateJWTToken(user.ID)
	if err != nil {
		return nil, "", time.Time{}, errors.New("failed to generate authentication token")
	}

	return &user, token, expiresAt, nil
}

// ExternalAuth authenticates or creates a user with external provider credentials
func (s *AuthService) ExternalAuth(dto *dtos.ExternalAuthDTO) (*entities.User, string, time.Time, bool, error) {
	// Validate auth type
	if dto.AuthType != "google" && dto.AuthType != "facebook" {
		return nil, "", time.Time{}, false, errors.New("invalid authentication type")
	}

	// TODO: Verify the external token with Google/Facebook API
	// This would normally involve sending the token to Google/Facebook
	// For now, we'll assume the token is valid

	db := databases.GetDB()

	// Start transaction
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// First try to find existing user by email
	var user entities.User
	var isNewUser bool = false
	if err := tx.Where("email = ?", dto.Email).First(&user).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			tx.Rollback()
			return nil, "", time.Time{}, false, err
		}

		// User doesn't exist, create a new one
		isNewUser = true
		email := dto.Email
		now := time.Now()
		user = entities.User{
			FirstName: dto.FirstName,
			LastName:  dto.LastName,
			Email:     &email,
			CreatedAt: now,
			UpdatedAt: now,
		}

		if err := tx.Create(&user).Error; err != nil {
			tx.Rollback()
			return nil, "", time.Time{}, false, err
		}
	}

	// Check if user already has credentials for this auth type
	var cred entities.UserAuthCredential
	if err := tx.Where("user_id = ? AND auth_type = ? AND external_id = ?", user.ID, dto.AuthType, dto.ExternalID).First(&cred).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			tx.Rollback()
			return nil, "", time.Time{}, false, err
		}

		// Credentials don't exist, create new ones
		now := time.Now()
		cred = entities.UserAuthCredential{
			UserID:     user.ID,
			AuthType:   dto.AuthType,
			ExternalID: &dto.ExternalID,
			CreatedAt:  now,
			UpdatedAt:  now,
		}

		if err := tx.Create(&cred).Error; err != nil {
			tx.Rollback()
			return nil, "", time.Time{}, false, err
		}
	}

	if err := tx.Commit().Error; err != nil {
		return nil, "", time.Time{}, false, err
	}

	// Generate JWT token
	token, expiresAt, err := helpers.GenerateJWTToken(user.ID)
	if err != nil {
		return nil, "", time.Time{}, false, errors.New("failed to generate authentication token")
	}

	return &user, token, expiresAt, isNewUser, nil
}

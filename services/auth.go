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

func (s *AuthService) RegisterUser(dto *dtos.RegisterDTO) (*entities.User, string, *time.Time, error) {
	db := databases.GetDB()
	var existingUser entities.User
	if err := db.Where("email = ?", dto.Email).First(&existingUser).Error; err == nil {
		return nil, "", nil, errors.New("user with this email already exists")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, "", nil, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(dto.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", nil, errors.New("failed to hash password")
	}

	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	email := dto.Email
	now := time.Now()
	user := entities.User{
		Username:  dto.Username,
		Email:     &email,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		return nil, "", nil, err
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
		return nil, "", nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, "", nil, err
	}

	token, expiresAt, err := helpers.GenerateJWTToken(user.ID)
	if err != nil {
		return nil, "", nil, errors.New("failed to generate authentication token")
	}

	return &user, token, expiresAt, nil
}

func (s *AuthService) LoginUser(email, password string) (*entities.User, string, *time.Time, error) {
	db := databases.GetDB()

	var user entities.User
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, "", nil, errors.New("invalid email or password")
	}

	var credentials entities.UserAuthCredential
	if err := db.Where("user_id = ? AND auth_type = ?", user.ID, "local").First(&credentials).Error; err != nil {
		return nil, "", nil, errors.New("invalid email or password")
	}

	if credentials.PasswordHash == nil {
		return nil, "", nil, errors.New("account has no password set")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(*credentials.PasswordHash), []byte(password)); err != nil {
		return nil, "", nil, errors.New("invalid email or password")
	}

	token, expiresAt, err := helpers.GenerateJWTToken(user.ID)
	if err != nil {
		return nil, "", nil, errors.New("failed to generate authentication token")
	}

	return &user, token, expiresAt, nil
}

func (s *AuthService) ExternalAuth(dto *dtos.ExternalAuthDTO) (*entities.User, string, *time.Time, bool, error) {
	if dto.AuthType != "google" && dto.AuthType != "facebook" {
		return nil, "", nil, false, errors.New("invalid authentication type")
	}

	db := databases.GetDB()

	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var user entities.User
	var isNewUser bool = false
	if err := tx.Where("email = ?", dto.Email).First(&user).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			tx.Rollback()
			return nil, "", nil, false, err
		}

		isNewUser = true
		email := dto.Email
		now := time.Now()
		user = entities.User{
			Username:  dto.Username,
			Email:     &email,
			CreatedAt: now,
			UpdatedAt: now,
		}

		if err := tx.Create(&user).Error; err != nil {
			tx.Rollback()
			return nil, "", nil, false, err
		}
	}

	var cred entities.UserAuthCredential
	if err := tx.Where("user_id = ? AND auth_type = ? AND external_id = ?", user.ID, dto.AuthType, dto.ExternalID).First(&cred).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			tx.Rollback()
			return nil, "", nil, false, err
		}

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
			return nil, "", nil, false, err
		}
	}

	if err := tx.Commit().Error; err != nil {
		return nil, "", nil, false, err
	}

	token, expiresAt, err := helpers.GenerateJWTToken(user.ID)
	if err != nil {
		return nil, "", nil, false, errors.New("failed to generate authentication token")
	}

	return &user, token, expiresAt, isNewUser, nil
}

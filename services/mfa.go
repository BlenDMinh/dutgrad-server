package services

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"image/png"
	"strconv"
	"strings"
	"time"

	"github.com/BlenDMinh/dutgrad-server/databases"
	"github.com/BlenDMinh/dutgrad-server/databases/entities"
	"github.com/BlenDMinh/dutgrad-server/databases/repositories"
	"github.com/BlenDMinh/dutgrad-server/helpers"
	"github.com/BlenDMinh/dutgrad-server/models/dtos"
	"github.com/pquerna/otp/totp"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type MFAService struct {
	redisService       *RedisService
	userRepo           *repositories.UserRepository
	userMFARepo        *repositories.UserMFARepository
	authCredentialRepo *repositories.UserAuthCredentialRepository
}

func NewMFAService() *MFAService {
	return &MFAService{
		redisService:       NewRedisService(),
		userRepo:           repositories.NewUserRepository(),
		userMFARepo:        repositories.NewUserMFARepository(),
		authCredentialRepo: repositories.NewUserAuthCredentialRepository(),
	}
}

func (s *MFAService) GenerateMFASetup(userID uint) (*dtos.MFASetupResponse, error) {
	user, err := s.userRepo.GetById(userID)
	if err != nil {
		return nil, err
	}

	existingMFA, err := s.userMFARepo.GetByUserID(userID)
	if err == nil {
		if existingMFA.Verified {
			return nil, errors.New("MFA is already enabled for this account")
		}
		s.userMFARepo.DeleteByUserID(userID)
	}

	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "DUTGrad",
		AccountName: *user.Email,
	})
	if err != nil {
		return nil, errors.New("failed to generate MFA secret")
	}

	backupCodes, err := s.generateBackupCodes(10)
	if err != nil {
		return nil, errors.New("failed to generate backup codes")
	}

	mfa := entities.UserMFA{
		UserID:      userID,
		Secret:      key.Secret(),
		BackupCodes: entities.BackupCodes(backupCodes),
		Verified:    false,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	createdMFA, err := s.userMFARepo.Create(&mfa)
	if err != nil {
		return nil, errors.New("failed to save MFA configuration")
	}

	qrCode, err := key.Image(200, 200)
	if err != nil {
		return nil, errors.New("failed to generate QR code")
	}

	var buf bytes.Buffer
	err = png.Encode(&buf, qrCode)
	if err != nil {
		return nil, errors.New("failed to encode QR code image")
	}

	dataURL := "data:image/png;base64," + base64.StdEncoding.EncodeToString(buf.Bytes())

	return &dtos.MFASetupResponse{
		Secret:          createdMFA.Secret,
		QRCodeDataURL:   dataURL,
		BackupCodes:     backupCodes,
		ProvisioningURI: key.URL(),
	}, nil
}

func (s *MFAService) VerifyMFASetup(userID uint, code string) error {
	mfa, err := s.userMFARepo.GetByUserID(userID)
	if err != nil {
		return errors.New("MFA setup not found")
	}

	valid := totp.Validate(code, mfa.Secret)
	if !valid {
		return errors.New("invalid verification code")
	}

	db := databases.GetDB()
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&entities.UserMFA{}).Where("id = ?", mfa.ID).Update("verified", true).Error; err != nil {
			return errors.New("failed to update MFA status")
		}

		if err := tx.Model(&entities.User{}).Where("id = ?", userID).Update("mfa_enabled", true).Error; err != nil {
			return errors.New("failed to enable MFA for user")
		}

		return nil
	})
}

func (s *MFAService) DisableMFA(userID uint) error {
	user, err := s.userRepo.GetById(userID)
	if err != nil {
		return err
	}

	if !user.MFAEnabled {
		return errors.New("MFA is not enabled for this account")
	}

	db := databases.GetDB()
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("user_id = ?", userID).Delete(&entities.UserMFA{}).Error; err != nil {
			return errors.New("failed to delete MFA configuration")
		}

		if err := tx.Model(&entities.User{}).Where("id = ?", userID).Update("mfa_enabled", false).Error; err != nil {
			return errors.New("failed to disable MFA for user")
		}

		return nil
	})
}

func (s *MFAService) VerifyMFACode(userID uint, code string, useBackupCode bool) bool {
	mfa, err := s.userMFARepo.GetByUserID(userID)
	if err != nil {
		return false
	}

	if !mfa.Verified {
		return false
	}

	if useBackupCode {
		for i, backupCode := range mfa.BackupCodes {
			if backupCode == code {
				newCodes := append(mfa.BackupCodes[:i], mfa.BackupCodes[i+1:]...)
				s.userMFARepo.UpdateBackupCodes(mfa.ID, newCodes)
				return true
			}
		}
		return false
	}

	return totp.Validate(code, mfa.Secret)
}

func (s *MFAService) CreateTempToken(userID uint) (string, time.Time, error) {
	token := make([]byte, 32)
	if _, err := rand.Read(token); err != nil {
		return "", time.Time{}, err
	}

	tempToken := base64.URLEncoding.EncodeToString(token)

	expiresAt := time.Now().Add(5 * time.Minute)

	key := fmt.Sprintf("mfa:temp:%s", tempToken)
	if err := s.redisService.Set(key, fmt.Sprintf("%d", userID), 5*time.Minute); err != nil {
		return "", time.Time{}, err
	}

	return tempToken, expiresAt, nil
}

func (s *MFAService) GetUserIDFromTempToken(token string) (uint, error) {
	key := fmt.Sprintf("mfa:temp:%s", token)
	value, err := s.redisService.Get(key)
	if err != nil {
		return 0, errors.New("invalid or expired token")
	}

	value = strings.Trim(value, "\"")

	userIDInt, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid token data: %v", err)
	}

	s.redisService.Del(key)

	return uint(userIDInt), nil
}

func (s *MFAService) generateBackupCodes(count int) ([]string, error) {
	codes := make([]string, count)

	for i := 0; i < count; i++ {
		randomBytes := make([]byte, 8)
		if _, err := rand.Read(randomBytes); err != nil {
			return nil, err
		}

		code := fmt.Sprintf("%02x%02x-%02x%02x",
			randomBytes[0], randomBytes[1], randomBytes[2], randomBytes[3])
		codes[i] = code
	}

	return codes, nil
}

func (s *MFAService) FirstFactorAuth(email, password string) (*entities.User, bool, error) {
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return nil, false, errors.New("invalid email or password")
	}

	credentials, err := s.authCredentialRepo.GetByUserIDAndType(user.ID, "local")
	if err != nil {
		return nil, false, errors.New("invalid email or password")
	}

	if credentials.PasswordHash == nil {
		return nil, false, errors.New("account has no password set")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(*credentials.PasswordHash), []byte(password)); err != nil {
		return nil, false, errors.New("invalid email or password")
	}

	if !user.MFAEnabled {
		return user, false, nil
	}

	return user, true, nil
}

func (s *MFAService) CompleteLogin(userID uint) (string, *time.Time, error) {
	token, expiresAt, err := helpers.GenerateJWTToken(userID)
	if err != nil {
		return "", nil, errors.New("failed to generate authentication token")
	}

	return token, expiresAt, nil
}

func (s *MFAService) GetUserMFAStatus(userID uint) (bool, error) {
	user, err := s.userRepo.GetById(userID)
	if err != nil {
		return false, err
	}

	return user.MFAEnabled, nil
}

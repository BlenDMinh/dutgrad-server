package repositories

import (
	"errors"
	"strings"
	"time"

	"github.com/BlenDMinh/dutgrad-server/databases"
	"github.com/BlenDMinh/dutgrad-server/databases/entities"
	"github.com/BlenDMinh/dutgrad-server/models/dtos"
)

// UserRepository handles database operations for users
type UserRepository struct {
	*CrudRepository[entities.User, uint]
}

// NewUserRepository creates a new user repository
func NewUserRepository() *UserRepository {
	return &UserRepository{
		CrudRepository: NewCrudRepository[entities.User, uint](),
	}
}

// GetSpacesByUserId retrieves spaces by user ID
func (r *UserRepository) GetSpacesByUserId(userId uint) ([]entities.Space, error) {
	var spaces []entities.Space
	db := databases.GetDB()
	err := db.Joins("JOIN space_users ON space_users.space_id = spaces.id").
		Where("space_users.user_id = ?", userId).
		Find(&spaces).Error
	return spaces, err
}

// GetByEmail gets a user by email
func (r *UserRepository) GetByEmail(email string) (*entities.User, error) {
	db := databases.GetDB()
	var user entities.User
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateMFAStatus updates the MFA status for a user
func (r *UserRepository) UpdateMFAStatus(userID uint, enabled bool) error {
	db := databases.GetDB()
	return db.Model(&entities.User{}).Where("id = ?", userID).Update("mfa_enabled", enabled).Error
}

// Transaction executes db operations in a transaction
func (r *UserRepository) Transaction(fn func(*databases.Transaction) error) error {
	db := databases.GetDB()
	tx := db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := fn(&databases.Transaction{DB: tx}); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// GetUserByEmail retrieves a user by email
func (r *UserRepository) GetUserByEmail(email string) (*entities.User, error) {
	var user entities.User
	db := databases.GetDB()
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetInvitationsByUserId retrieves invitations by user ID
func (r *UserRepository) GetInvitationsByUserId(InvitedUserId uint) ([]entities.SpaceInvitation, error) {
	var invitations []entities.SpaceInvitation
	db := databases.GetDB()
	err := db.Preload("Space").Preload("Inviter").
		Where("invited_user_id = ?", InvitedUserId).Find(&invitations).Error
	if err != nil {
		return nil, err
	}
	return invitations, nil
}

// SearchUsers searches users by query
func (r *UserRepository) SearchUsers(query string) ([]entities.User, error) {
	var users []entities.User
	db := databases.GetDB()

	if strings.Contains(query, "@") {
		err := db.Where("email = ?", query).Find(&users).Error
		if err != nil {
			return nil, err
		}
	} else {
		err := db.Where("username ILIKE ?", "%"+query+"%").Find(&users).Error
		if err != nil {
			return nil, err
		}
	}

	return users, nil
}

// GetUserTier retrieves the tier of a user
func (r *UserRepository) GetUserTier(userID uint) (*entities.Tier, error) {
	db := databases.GetDB()
	var user entities.User

	if err := db.Preload("Tier").First(&user, userID).Error; err != nil {
		return nil, err
	}

	if user.TierID == nil || user.Tier == nil {
		return nil, errors.New("user has no tier information")
	}

	return user.Tier, nil
}

// GetUserTierUsage retrieves the tier usage of a user
func (s *UserRepository) GetUserTierUsage(userID uint) (*dtos.TierUsageResponse, error) {
	tier, err := s.GetUserTier(userID)
	if err != nil {
		return nil, err
	}

	db := databases.GetDB()
	var response dtos.TierUsageResponse
	response.Tier = tier
	response.Usage = &dtos.TierUsage{}

	err = db.Model(&entities.SpaceUser{}).
		Where("user_id = ? AND space_role_id = ?", userID, entities.Owner).
		Count(&response.Usage.SpaceCount).Error
	if err != nil {
		return nil, err
	}

	err = db.Model(&entities.Document{}).
		Joins("JOIN space_users ON space_users.space_id = documents.space_id").
		Where("space_users.user_id = ? AND space_users.space_role_id = ?", userID, entities.Owner).
		Count(&response.Usage.DocumentCount).Error
	if err != nil {
		return nil, err
	}

	err = db.Model(&entities.ChatHistory{}).
		Joins("JOIN user_query_sessions ON user_query_sessions.id = chat_histories.session_id").
		Where("user_query_sessions.user_id = ?", userID).
		Count(&response.Usage.QueryHistoryCount).Error
	if err != nil {
		return nil, err
	}

	today := time.Now().Format("2006-01-02")
	err = db.Model(&entities.ChatHistory{}).
		Joins("JOIN user_query_sessions ON user_query_sessions.id = chat_histories.session_id").
		Where("user_query_sessions.user_id = ? AND DATE(chat_histories.created_at) = ?", userID, today).
		Count(&response.Usage.TodayQueryCount).Error
	if err != nil {
		return nil, err
	}

	response.Usage.TodayApiCallCount = 0

	return &response, nil
}

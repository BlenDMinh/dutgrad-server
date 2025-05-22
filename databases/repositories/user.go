package repositories

import (
	"errors"
	"strings"
	"time"

	"github.com/BlenDMinh/dutgrad-server/databases"
	"github.com/BlenDMinh/dutgrad-server/databases/entities"
	"github.com/BlenDMinh/dutgrad-server/models/dtos"
)

type UserRepository interface {
	ICrudRepository[entities.User, uint]
	GetSpacesByUserId(userId uint) ([]entities.Space, error)
	GetByEmail(email string) (*entities.User, error)
	UpdateMFAStatus(userID uint, enabled bool) error
	GetInvitationsByUserId(InvitedUserId uint) ([]entities.SpaceInvitation, error)
	SearchUsers(query string) ([]entities.User, error)
	GetUserTier(userID uint) (*entities.Tier, error)
	GetUserTierUsage(userID uint) (*dtos.TierUsageResponse, error)
}

type UserRepositoryImpl struct {
	*CrudRepository[entities.User, uint]
}

func NewUserRepository() UserRepository {
	return &UserRepositoryImpl{
		CrudRepository: NewCrudRepository[entities.User, uint](),
	}
}

func (r *UserRepositoryImpl) GetSpacesByUserId(userId uint) ([]entities.Space, error) {
	db := databases.GetDB()
	var spaceUsers []entities.SpaceUser
	err := db.Preload("Space").Preload("SpaceRole").
		Where("user_id = ?", userId).
		Find(&spaceUsers).Error

	if err != nil {
		return nil, err
	}
	var userSpaces []entities.Space
	for _, spaceUser := range spaceUsers {
		userSpaces = append(userSpaces, spaceUser.Space)
	}

	return userSpaces, nil
}

func (r *UserRepositoryImpl) GetByEmail(email string) (*entities.User, error) {
	db := databases.GetDB()
	var user entities.User
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepositoryImpl) UpdateMFAStatus(userID uint, enabled bool) error {
	db := databases.GetDB()
	return db.Model(&entities.User{}).Where("id = ?", userID).Update("mfa_enabled", enabled).Error
}

func (r *UserRepositoryImpl) Transaction(fn func(*databases.Transaction) error) error {
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

func (r *UserRepositoryImpl) GetInvitationsByUserId(InvitedUserId uint) ([]entities.SpaceInvitation, error) {
	var invitations []entities.SpaceInvitation
	db := databases.GetDB()
	err := db.Preload("Space").Preload("Inviter").
		Where("invited_user_id = ?", InvitedUserId).Find(&invitations).Error
	if err != nil {
		return nil, err
	}
	return invitations, nil
}

func (r *UserRepositoryImpl) SearchUsers(query string) ([]entities.User, error) {
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

func (r *UserRepositoryImpl) GetUserTier(userID uint) (*entities.Tier, error) {
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

func (s *UserRepositoryImpl) GetUserTierUsage(userID uint) (*dtos.TierUsageResponse, error) {
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

package repositories

import (
	"github.com/BlenDMinh/dutgrad-server/databases"
	"github.com/BlenDMinh/dutgrad-server/databases/entities"
)

type UserQuerySessionRepository struct {
	*CrudRepository[entities.UserQuerySession, uint]
}

func NewUserQuerySessionRepository() *UserQuerySessionRepository {
	return &UserQuerySessionRepository{
		CrudRepository: NewCrudRepository[entities.UserQuerySession, uint](),
	}
}

func (s *UserQuerySessionRepository) CountByUserID(userID uint) (int64, error) {
	var count int64
	db := databases.GetDB()
	err := db.Model(&entities.UserQuerySession{}).Where("user_id = ?", userID).Count(&count).Error
	return count, err
}

func (s *UserQuerySessionRepository) GetByUserID(userID uint) ([]entities.UserQuerySession, error) {
	var sessions []entities.UserQuerySession
	db := databases.GetDB()
	err := db.Where("user_id = ?", userID).Find(&sessions).Error
	return sessions, err
}

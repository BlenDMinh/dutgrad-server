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
	err := db.Where("user_query_sessions.user_id = ?", userID).
		Joins("INNER JOIN space_users ON space_users.space_id = user_query_sessions.space_id AND space_users.user_id = ?", userID).
		Joins("LEFT JOIN chat_histories ON chat_histories.session_id = user_query_sessions.id").
		Group("user_query_sessions.id").
		Having("COUNT(chat_histories.id) > 0").
		Order("MAX(chat_histories.created_at) DESC").
		Preload("ChatHistories").
		Find(&sessions).
		Error
	return sessions, err
}

func (s *UserQuerySessionRepository) GetTempMessageByID(id uint) (*string, error) {
	var session entities.UserQuerySession
	db := databases.GetDB()
	err := db.Select("temp_message").Where("id = ?", id).First(&session).Error

	if err != nil {
		return nil, err
	}
	return session.TempMessage, nil
}

package repositories

import (
	"encoding/json"
	"fmt"

	"github.com/BlenDMinh/dutgrad-server/databases"
	"github.com/BlenDMinh/dutgrad-server/databases/entities"
)

type UserQuerySessionRepository interface {
	ICrudRepository[entities.UserQuerySession, uint]
	CountByUserID(userID uint) (int64, error)
	GetByUserID(userID uint) ([]entities.UserQuerySession, error)
	GetTempMessageByID(id uint) (*string, error)
	GetChatHistoryBySessionID(sessionID uint) ([]map[string]interface{}, error)
	ClearChatHistory(sessionID uint) error
}

type userQuerySessionRepositoryImpl struct {
	*CrudRepository[entities.UserQuerySession, uint]
}

func NewUserQuerySessionRepository() UserQuerySessionRepository {
	return &userQuerySessionRepositoryImpl{
		CrudRepository: NewCrudRepository[entities.UserQuerySession, uint](),
	}
}

func (s *userQuerySessionRepositoryImpl) CountByUserID(userID uint) (int64, error) {
	var count int64
	db := databases.GetDB()
	err := db.Model(&entities.UserQuerySession{}).Where("user_id = ?", userID).Count(&count).Error
	return count, err
}

func (s *userQuerySessionRepositoryImpl) GetByUserID(userID uint) ([]entities.UserQuerySession, error) {
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

func (s *userQuerySessionRepositoryImpl) GetTempMessageByID(id uint) (*string, error) {
	var session entities.UserQuerySession
	db := databases.GetDB()
	err := db.Select("temp_message").Where("id = ?", id).First(&session).Error

	if err != nil {
		return nil, err
	}
	return session.TempMessage, nil
}

func (s *userQuerySessionRepositoryImpl) GetChatHistoryBySessionID(sessionID uint) ([]map[string]interface{}, error) {
	var chatHistories []entities.ChatHistory
	db := databases.GetDB()

	err := db.Where("session_id = ?", sessionID).
		Order("created_at ASC").
		Find(&chatHistories).
		Error

	if err != nil {
		return nil, err
	}
	var result []map[string]interface{}
	for _, history := range chatHistories {
		var dbMessage map[string]interface{}
		if err := json.Unmarshal(history.Message, &dbMessage); err != nil {
			return nil, fmt.Errorf("failed to unmarshal message: %v", err)
		}

		messageType, _ := dbMessage["type"].(string)
		content, _ := dbMessage["content"].(string)

		message := map[string]interface{}{
			"id":        fmt.Sprintf("%d", history.ID),
			"content":   content,
			"isUser":    messageType == "human",
			"timestamp": history.CreatedAt,
		}

		result = append(result, message)
	}

	return result, nil
}

func (s *userQuerySessionRepositoryImpl) ClearChatHistory(sessionID uint) error {
	db := databases.GetDB()

	err := db.Where("session_id = ?", sessionID).Delete(&entities.ChatHistory{}).Error
	if err != nil {
		return fmt.Errorf("failed to clear chat history: %v", err)
	}

	err = db.Where("query_session_id = ?", sessionID).Delete(&entities.UserQuery{}).Error
	if err != nil {
		return fmt.Errorf("failed to delete user queries: %v", err)
	}

	err = db.Where("id = ?", sessionID).Delete(&entities.UserQuerySession{}).Error
	if err != nil {
		return fmt.Errorf("failed to delete session: %v", err)
	}

	return nil
}

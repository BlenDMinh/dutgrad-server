package services

import (
	"fmt"

	"github.com/BlenDMinh/dutgrad-server/databases/entities"
	"github.com/BlenDMinh/dutgrad-server/databases/repositories"
)

type UserQuerySessionService struct {
	CrudService[entities.UserQuerySession, uint]
}

func NewUserQuerySessionService() *UserQuerySessionService {
	return &UserQuerySessionService{
		CrudService: *NewCrudService(repositories.NewUserQuerySessionRepository()),
	}
}

func (s *UserQuerySessionService) GetChatSessionsByUserID(userID uint) ([]entities.UserQuerySession, error) {
	return s.repo.(*repositories.UserQuerySessionRepository).GetByUserID(userID)

}

func (s *UserQuerySessionService) CountChatSessionsByUserID(userID uint) (int64, error) {
	return s.repo.(*repositories.UserQuerySessionRepository).CountByUserID(userID)
}

func (s *UserQuerySessionService) GetTempMessageByID(id uint) (*string, error) {
	return s.repo.(*repositories.UserQuerySessionRepository).GetTempMessageByID(id)
}

func (s *UserQuerySessionService) GetChatHistoryBySessionID(sessionID uint, userID uint) ([]map[string]interface{}, error) {
	return s.repo.(*repositories.UserQuerySessionRepository).GetChatHistoryBySessionID(sessionID)
}

func (s *UserQuerySessionService) ClearChatHistoryBySessionID(sessionID uint, userID uint) error {
	session, err := s.GetById(sessionID)
	if err != nil {
		return fmt.Errorf("failed to find session: %v", err)
	}

	if session.UserID == nil || *session.UserID != userID {
		return fmt.Errorf("unauthorized: you can only clear chat history for your own sessions")
	}

	return s.repo.(*repositories.UserQuerySessionRepository).ClearChatHistory(sessionID)
}

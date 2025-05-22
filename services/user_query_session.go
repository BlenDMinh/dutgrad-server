package services

import (
	"fmt"

	"github.com/BlenDMinh/dutgrad-server/databases/entities"
	"github.com/BlenDMinh/dutgrad-server/databases/repositories"
)

type UserQuerySessionService interface {
	ICrudService[entities.UserQuerySession, uint]
	GetChatSessionsByUserID(userID uint) ([]entities.UserQuerySession, error)
	CountChatSessionsByUserID(userID uint) (int64, error)
	GetTempMessageByID(id uint) (*string, error)
	GetChatHistoryBySessionID(sessionID uint, userID uint) ([]map[string]interface{}, error)
	ClearChatHistoryBySessionID(sessionID uint, userID uint) error
}

type UserQuerySessionServiceImpl struct {
	CrudService[entities.UserQuerySession, uint]
	repo repositories.UserQuerySessionRepository
}

func NewUserQuerySessionService() UserQuerySessionService {
	crudService := NewCrudService(repositories.NewUserQuerySessionRepository())
	repo := crudService.repo.(repositories.UserQuerySessionRepository)
	return &UserQuerySessionServiceImpl{
		CrudService: *crudService,
		repo:        repo,
	}
}

func (s *UserQuerySessionServiceImpl) GetChatSessionsByUserID(userID uint) ([]entities.UserQuerySession, error) {
	return s.repo.GetByUserID(userID)

}

func (s *UserQuerySessionServiceImpl) CountChatSessionsByUserID(userID uint) (int64, error) {
	return s.repo.CountByUserID(userID)
}

func (s *UserQuerySessionServiceImpl) GetTempMessageByID(id uint) (*string, error) {
	return s.repo.GetTempMessageByID(id)
}

func (s *UserQuerySessionServiceImpl) GetChatHistoryBySessionID(sessionID uint, userID uint) ([]map[string]interface{}, error) {
	return s.repo.GetChatHistoryBySessionID(sessionID)
}

func (s *UserQuerySessionServiceImpl) ClearChatHistoryBySessionID(sessionID uint, userID uint) error {
	session, err := s.GetById(sessionID)
	if err != nil {
		return fmt.Errorf("failed to find session: %v", err)
	}

	if session.UserID == nil || *session.UserID != userID {
		return fmt.Errorf("unauthorized: you can only clear chat history for your own sessions")
	}

	return s.repo.ClearChatHistory(sessionID)
}

package services

import (
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

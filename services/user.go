package services

import (
	"github.com/BlenDMinh/dutgrad-server/databases/entities"
	"github.com/BlenDMinh/dutgrad-server/databases/repositories"
)

type UserService struct {
	CrudService[entities.User, uint]
}

func NewUserService() *UserService {
	return &UserService{
		CrudService: *NewCrudService(repositories.NewUserRepository()),
	}
}

func (s *UserService) GetSpacesByUserId(userId uint) ([]entities.Space, error) {
	return s.repo.(*repositories.UserRepository).GetSpacesByUserId(userId)
}

func (s *UserService) GetUserByEmail(email string) (*entities.User, error) {
	return s.repo.(*repositories.UserRepository).GetUserByEmail(email)
}

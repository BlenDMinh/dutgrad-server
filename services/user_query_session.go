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

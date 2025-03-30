package services

import (
	"github.com/BlenDMinh/dutgrad-server/databases/entities"
	"github.com/BlenDMinh/dutgrad-server/databases/repositories"
)

type UserQueryService struct {
	CrudService[entities.UserQuery, uint]
}

func NewUserQueryService() *UserQueryService {
	return &UserQueryService{
		CrudService: *NewCrudService(repositories.NewUserQueryRepository()),
	}
}

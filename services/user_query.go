package services

import (
	"github.com/BlenDMinh/dutgrad-server/databases/entities"
	"github.com/BlenDMinh/dutgrad-server/databases/repositories"
)

type UserQueryService interface {
	ICrudService[entities.UserQuery, uint]
}

type UserQueryServiceImpl struct {
	CrudService[entities.UserQuery, uint]
	repo repositories.UserQueryRepository
}

func NewUserQueryService() UserQueryService {
	crudService := NewCrudService(repositories.NewUserQueryRepository())
	repo := crudService.repo.(repositories.UserQueryRepository)
	return &UserQueryServiceImpl{
		CrudService: *crudService,
		repo:        repo,
	}
}

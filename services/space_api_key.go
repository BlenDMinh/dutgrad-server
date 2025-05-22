package services

import (
	"github.com/BlenDMinh/dutgrad-server/databases/entities"
	"github.com/BlenDMinh/dutgrad-server/databases/repositories"
)

type SpaceApiKeyService interface {
	ICrudService[entities.SpaceAPIKey, uint]
	GetAllBySpaceID(spaceID uint) ([]entities.SpaceAPIKey, error)
}

type spaceApiKeyServiceImpl struct {
	CrudService[entities.SpaceAPIKey, uint]
	repo repositories.SpaceApiKeyRepository
}

func NewSpaceApiKeyService() SpaceApiKeyService {
	crudService := NewCrudService(repositories.NewSpaceApiKeyRepository())
	repo := crudService.repo.(repositories.SpaceApiKeyRepository)
	return &spaceApiKeyServiceImpl{
		CrudService: *crudService,
		repo:        repo,
	}
}

func (s *spaceApiKeyServiceImpl) GetAllBySpaceID(spaceID uint) ([]entities.SpaceAPIKey, error) {
	return s.repo.GetAllBySpaceID(spaceID)
}

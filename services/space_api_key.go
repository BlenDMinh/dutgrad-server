package services

import (
	"github.com/BlenDMinh/dutgrad-server/databases/entities"
	"github.com/BlenDMinh/dutgrad-server/databases/repositories"
)

type SpaceApiKeyService struct {
	CrudService[entities.SpaceAPIKey, uint]
}

func NewSpaceApiKeyService() *SpaceApiKeyService {
	return &SpaceApiKeyService{
		CrudService: *NewCrudService(repositories.NewSpaceApiKeyRepository()),
	}
}

func (s *SpaceApiKeyService) GetAllBySpaceID(spaceID uint) ([]entities.SpaceAPIKey, error) {
	return s.repo.(*repositories.SpaceApiKeyRepository).GetAllBySpaceID(spaceID)
}

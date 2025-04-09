package services

import (
	"github.com/BlenDMinh/dutgrad-server/databases/entities"
	"github.com/BlenDMinh/dutgrad-server/databases/repositories"
)

type SpaceService struct {
	CrudService[entities.Space, uint]
}

func NewSpaceService() *SpaceService {
	return &SpaceService{
		CrudService: *NewCrudService(repositories.NewSpaceRepository()),
	}
}

func (s *SpaceService) GetPublicSpaces() ([]entities.Space, error) {
	return s.repo.(*repositories.SpaceRepository).FindPublicSpaces()
}

func (s *SpaceService) GetMembers(spaceId uint) ([]entities.SpaceUser, error) {
	return s.repo.(*repositories.SpaceRepository).GetMembers(spaceId)
}

func (s *SpaceService) GetInvitations(spaceId uint) ([]entities.SpaceInvitation, error) {
	return s.repo.(*repositories.SpaceRepository).GetInvitations(spaceId)
}

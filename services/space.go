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

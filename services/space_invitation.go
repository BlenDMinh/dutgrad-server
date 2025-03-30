package services

import (
	"github.com/BlenDMinh/dutgrad-server/databases/entities"
	"github.com/BlenDMinh/dutgrad-server/databases/repositories"
)

type SpaceInvitationService struct {
	CrudService[entities.SpaceInvitation, uint]
}

func NewSpaceInvitationService() *SpaceInvitationService {
	return &SpaceInvitationService{
		CrudService: *NewCrudService(repositories.NewSpaceInvitationRepository()),
	}
}

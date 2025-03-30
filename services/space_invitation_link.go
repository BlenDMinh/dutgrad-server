package services

import (
	"github.com/BlenDMinh/dutgrad-server/databases/entities"
	"github.com/BlenDMinh/dutgrad-server/databases/repositories"
)

type SpaceInvitationLinkService struct {
	CrudService[entities.SpaceInvitationLink, uint]
}

func NewSpaceInvitationLinkService() *SpaceInvitationLinkService {
	return &SpaceInvitationLinkService{
		CrudService: *NewCrudService(repositories.NewSpaceInvitationLinkRepository()),
	}
}

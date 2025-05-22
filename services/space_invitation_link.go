package services

import (
	"github.com/BlenDMinh/dutgrad-server/databases/entities"
	"github.com/BlenDMinh/dutgrad-server/databases/repositories"
)

type SpaceInvitationLinkService interface {
	ICrudService[entities.SpaceInvitationLink, uint]
}

type spaceInvitationLinkServiceImpl struct {
	CrudService[entities.SpaceInvitationLink, uint]
	repo repositories.SpaceInvitationLinkRepository
}

func NewSpaceInvitationLinkService() SpaceInvitationLinkService {
	crudService := NewCrudService(repositories.NewSpaceInvitationLinkRepository())
	repo := crudService.repo.(repositories.SpaceInvitationLinkRepository)
	return &spaceInvitationLinkServiceImpl{
		CrudService: *crudService,
		repo:        repo,
	}
}

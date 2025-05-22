package controllers

import (
	"github.com/BlenDMinh/dutgrad-server/databases/entities"
	"github.com/BlenDMinh/dutgrad-server/services"
)

type SpaceInvitationLinkController struct {
	CrudController[entities.SpaceInvitationLink, uint]
	service services.SpaceInvitationLinkService
}

func NewSpaceInvitationLinkController(
	service services.SpaceInvitationLinkService,
) *SpaceInvitationLinkController {
	crudController := NewCrudController(service)
	return &SpaceInvitationLinkController{
		CrudController: *crudController,
		service:        service,
	}
}

package controllers

import (
	"github.com/BlenDMinh/dutgrad-server/databases/entities"
	"github.com/BlenDMinh/dutgrad-server/services"
)

type SpaceInvitationLinkController struct {
	CrudController[entities.SpaceInvitationLink, uint]
}

func NewSpaceInvitationLinkController() *SpaceInvitationLinkController {
	return &SpaceInvitationLinkController{
		CrudController: *NewCrudController(services.NewSpaceInvitationLinkService()),
	}
}

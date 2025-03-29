package controllers

import (
	"github.com/BlenDMinh/dutgrad-server/databases/entities"
	"github.com/BlenDMinh/dutgrad-server/services"
)

type SpaceInvitationController struct {
	CrudController[entities.SpaceInvitation, uint]
}

func NewSpaceInvitationController() *SpaceInvitationController {
	return &SpaceInvitationController{
		CrudController: *NewCrudController(services.NewSpaceInvitationService()),
	}
}

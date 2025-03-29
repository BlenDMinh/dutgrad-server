package controllers

import (
	"github.com/BlenDMinh/dutgrad-server/databases/entities"
	"github.com/BlenDMinh/dutgrad-server/services"
)

type UserQuerySessionController struct {
	CrudController[entities.UserQuerySession, uint]
}

func NewUserQuerySessionController() *UserQuerySessionController {
	return &UserQuerySessionController{
		CrudController: *NewCrudController(services.NewUserQuerySessionService()),
	}
}

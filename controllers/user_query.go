package controllers

import (
	"github.com/BlenDMinh/dutgrad-server/databases/entities"
	"github.com/BlenDMinh/dutgrad-server/services"
)

type UserQueryController struct {
	CrudController[entities.UserQuery, uint]
}

func NewUserQueryController() *UserQueryController {
	return &UserQueryController{
		CrudController: *NewCrudController(services.NewUserQueryService()),
	}
}

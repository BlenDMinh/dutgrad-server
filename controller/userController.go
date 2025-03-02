package controller

import (
	"github.com/BlenDMinh/dutgrad-server/database/entity"
)

type UserController struct {
	CrudController
}

func (c *UserController) getModel() interface{} {
	return &entity.User{}
}

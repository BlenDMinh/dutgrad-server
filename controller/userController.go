package controller

import (
	"github.com/BlenDMinh/dutgrad-server/database/entity"
)

type UserController struct {
	CrudController[entity.User]
}

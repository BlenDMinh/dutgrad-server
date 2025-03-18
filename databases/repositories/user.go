package repositories

import "github.com/BlenDMinh/dutgrad-server/databases/entities"

type UserRepository struct {
	*CrudRepository[*entities.User, uint]
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		CrudRepository: NewCrudRepository[*entities.User, uint](),
	}
}

package repositories

import "github.com/BlenDMinh/dutgrad-server/databases/entities"

type UserQuerySessionRepository struct {
	*CrudRepository[entities.UserQuerySession, uint]
}

func NewUserQuerySessionRepository() *UserQuerySessionRepository {
	return &UserQuerySessionRepository{
		CrudRepository: NewCrudRepository[entities.UserQuerySession, uint](),
	}
}

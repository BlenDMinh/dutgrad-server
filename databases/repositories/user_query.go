package repositories

import "github.com/BlenDMinh/dutgrad-server/databases/entities"

type UserQueryRepository struct {
	*CrudRepository[entities.UserQuery, uint]
}

func NewUserQueryRepository() *UserQueryRepository {
	return &UserQueryRepository{
		CrudRepository: NewCrudRepository[entities.UserQuery, uint](),
	}
}

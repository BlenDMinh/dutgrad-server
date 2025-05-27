package repositories

import "github.com/BlenDMinh/dutgrad-server/databases/entities"

type UserQueryRepository interface {
	ICrudRepository[entities.UserQuery, uint]
}

type userQueryRepositoryImpl struct {
	*CrudRepository[entities.UserQuery, uint]
}

func NewUserQueryRepository() UserQueryRepository {
	return &userQueryRepositoryImpl{
		CrudRepository: NewCrudRepository[entities.UserQuery, uint](),
	}
}

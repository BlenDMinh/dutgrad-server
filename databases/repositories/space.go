package repositories

import (
	"github.com/BlenDMinh/dutgrad-server/databases/entities"
)

type SpaceRepository struct {
	*CrudRepository[entities.Space, uint]
}

func NewSpaceRepository() *SpaceRepository {
	return &SpaceRepository{
		CrudRepository: NewCrudRepository[entities.Space, uint](),
	}
}

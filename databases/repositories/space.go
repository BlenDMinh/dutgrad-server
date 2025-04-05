package repositories

import (
	"github.com/BlenDMinh/dutgrad-server/databases"
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

func (r *SpaceRepository) FindPublicSpaces() ([]entities.Space, error) {
	var spaces []entities.Space
	db := databases.GetDB()
	err := db.Where("privacy_status = ?", false).Find(&spaces).Error
	return spaces, err
}

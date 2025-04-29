package repositories

import (
	"github.com/BlenDMinh/dutgrad-server/databases"
	"github.com/BlenDMinh/dutgrad-server/databases/entities"
)

type SpaceApiKeyRepository struct {
	*CrudRepository[entities.SpaceAPIKey, uint]
}

func NewSpaceApiKeyRepository() *SpaceApiKeyRepository {
	return &SpaceApiKeyRepository{
		CrudRepository: NewCrudRepository[entities.SpaceAPIKey, uint](),
	}
}

func (s *SpaceApiKeyRepository) GetAllBySpaceID(spaceID uint) ([]entities.SpaceAPIKey, error) {
	var results []entities.SpaceAPIKey
	db := databases.GetDB()
	err := db.Where("space_id = ?", spaceID).Find(&results).Error
	return results, err
}

package repositories

import (
	"github.com/BlenDMinh/dutgrad-server/databases"
	"github.com/BlenDMinh/dutgrad-server/databases/entities"
)

type DocumentRepository struct {
	*CrudRepository[entities.Document, uint]
}

func NewDocumentRepository() *DocumentRepository {
	return &DocumentRepository{
		CrudRepository: NewCrudRepository[entities.Document, uint](),
	}
}

func (r *DocumentRepository) GetBySpaceID(spaceID uint) ([]entities.Document, error) {
	db := databases.GetDB()
	documents := []entities.Document{}
	err := db.Where("space_id = ?", spaceID).Find(&documents).Error
	if err != nil {
		return nil, err
	}
	return documents, nil
}

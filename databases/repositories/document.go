package repositories

import (
	"github.com/BlenDMinh/dutgrad-server/databases"
	"github.com/BlenDMinh/dutgrad-server/databases/entities"
)

type DocumentRepository interface {
	ICrudRepository[entities.Document, uint]
	GetBySpaceID(spaceID uint) ([]entities.Document, error)
	CountUserDocuments(userID uint) (int64, error)
}

type documentRepositoryImpl struct {
	*CrudRepository[entities.Document, uint]
}

func NewDocumentRepository() DocumentRepository {
	return &documentRepositoryImpl{
		CrudRepository: NewCrudRepository[entities.Document, uint](),
	}
}

func (r *documentRepositoryImpl) GetBySpaceID(spaceID uint) ([]entities.Document, error) {
	db := databases.GetDB()
	documents := []entities.Document{}
	err := db.Where("space_id = ?", spaceID).Find(&documents).Error
	if err != nil {
		return nil, err
	}
	return documents, nil
}

func (s *documentRepositoryImpl) CountUserDocuments(userID uint) (int64, error) {
	var count int64
	db := databases.GetDB()

	err := db.Model(&entities.Document{}).
		Joins("JOIN spaces ON documents.space_id = spaces.id").
		Joins("JOIN space_users ON spaces.id = space_users.space_id").
		Where("space_users.user_id = ?", userID).
		Count(&count).Error

	return count, err
}

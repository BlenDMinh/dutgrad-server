package repositories

import (
	"errors"

	"github.com/BlenDMinh/dutgrad-server/databases"
	"github.com/BlenDMinh/dutgrad-server/databases/entities"
)

type DocumentRepository interface {
	ICrudRepository[entities.Document, uint]
	GetBySpaceID(spaceID uint) ([]entities.Document, error)
	GetUserRoleInSpace(userID, spaceID uint) (string, error)
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

func (s *documentRepositoryImpl) GetUserRoleInSpace(userID, spaceID uint) (string, error) {
	var spaceUser entities.SpaceUser
	db := databases.GetDB()
	result := db.Preload("SpaceRole").Where("user_id = ? AND space_id = ?", userID, spaceID).First(&spaceUser)
	if result.Error != nil {
		return "", result.Error
	}
	if spaceUser.SpaceRole.Name == "" {
		return "", errors.New("user has no role in this space")
	}
	return spaceUser.SpaceRole.Name, nil
}

package repositories

import (
	"errors"

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

func (s *DocumentRepository) GetUserRoleInSpace(userID, spaceID uint) (string, error) {
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

func (s *DocumentRepository) DeleteDocumentByID(documentID uint) error {
	return databases.GetDB().Delete(&entities.Document{}, documentID).Error
}

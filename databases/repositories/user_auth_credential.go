package repositories

import (
	"github.com/BlenDMinh/dutgrad-server/databases"
	"github.com/BlenDMinh/dutgrad-server/databases/entities"
)

// UserAuthCredentialRepository handles database operations for user authentication credentials
type UserAuthCredentialRepository struct {
	*CrudRepository[entities.UserAuthCredential, uint]
}

// NewUserAuthCredentialRepository creates a new UserAuthCredential repository
func NewUserAuthCredentialRepository() *UserAuthCredentialRepository {
	return &UserAuthCredentialRepository{
		CrudRepository: NewCrudRepository[entities.UserAuthCredential, uint](),
	}
}

// GetByUserIDAndType gets user auth credentials by user ID and auth type
func (r *UserAuthCredentialRepository) GetByUserIDAndType(userID uint, authType string) (*entities.UserAuthCredential, error) {
	db := databases.GetDB()
	var credential entities.UserAuthCredential
	if err := db.Where("user_id = ? AND auth_type = ?", userID, authType).First(&credential).Error; err != nil {
		return nil, err
	}
	return &credential, nil
}

// GetByExternalIDAndType gets user auth credentials by external ID and auth type
func (r *UserAuthCredentialRepository) GetByExternalIDAndType(externalID string, authType string) (*entities.UserAuthCredential, error) {
	db := databases.GetDB()
	var credential entities.UserAuthCredential
	if err := db.Where("external_id = ? AND auth_type = ?", externalID, authType).First(&credential).Error; err != nil {
		return nil, err
	}
	return &credential, nil
}

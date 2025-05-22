package repositories

import (
	"github.com/BlenDMinh/dutgrad-server/databases"
	"github.com/BlenDMinh/dutgrad-server/databases/entities"
)

type UserMFARepository interface {
	ICrudRepository[entities.UserMFA, uint]
	DeleteByUserID(userID uint) error
	GetByUserID(userID uint) (*entities.UserMFA, error)
	UpdateBackupCodes(mfaID uint, backupCodes entities.BackupCodes) error
	UpdateVerificationStatus(mfaID uint, verified bool) error
}

type userMFARepositoryImpl struct {
	*CrudRepository[entities.UserMFA, uint]
}

func NewUserMFARepository() UserMFARepository {
	return &userMFARepositoryImpl{
		CrudRepository: NewCrudRepository[entities.UserMFA, uint](),
	}
}

func (r *userMFARepositoryImpl) GetByUserID(userID uint) (*entities.UserMFA, error) {
	db := databases.GetDB()
	var mfa entities.UserMFA
	if err := db.Where("user_id = ?", userID).First(&mfa).Error; err != nil {
		return nil, err
	}
	return &mfa, nil
}

func (r *userMFARepositoryImpl) DeleteByUserID(userID uint) error {
	db := databases.GetDB()
	return db.Where("user_id = ?", userID).Delete(&entities.UserMFA{}).Error
}

func (r *userMFARepositoryImpl) UpdateBackupCodes(mfaID uint, backupCodes entities.BackupCodes) error {
	db := databases.GetDB()
	return db.Model(&entities.UserMFA{}).Where("id = ?", mfaID).Update("backup_codes", backupCodes).Error
}

func (r *userMFARepositoryImpl) UpdateVerificationStatus(mfaID uint, verified bool) error {
	db := databases.GetDB()
	return db.Model(&entities.UserMFA{}).Where("id = ?", mfaID).Update("verified", verified).Error
}

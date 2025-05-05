package repositories

import (
	"github.com/BlenDMinh/dutgrad-server/databases"
	"github.com/BlenDMinh/dutgrad-server/databases/entities"
)

type UserMFARepository struct {
	*CrudRepository[entities.UserMFA, uint]
}

func NewUserMFARepository() *UserMFARepository {
	return &UserMFARepository{
		CrudRepository: NewCrudRepository[entities.UserMFA, uint](),
	}
}

func (r *UserMFARepository) GetByUserID(userID uint) (*entities.UserMFA, error) {
	db := databases.GetDB()
	var mfa entities.UserMFA
	if err := db.Where("user_id = ?", userID).First(&mfa).Error; err != nil {
		return nil, err
	}
	return &mfa, nil
}

func (r *UserMFARepository) DeleteByUserID(userID uint) error {
	db := databases.GetDB()
	return db.Where("user_id = ?", userID).Delete(&entities.UserMFA{}).Error
}

func (r *UserMFARepository) UpdateBackupCodes(mfaID uint, backupCodes entities.BackupCodes) error {
	db := databases.GetDB()
	return db.Model(&entities.UserMFA{}).Where("id = ?", mfaID).Update("backup_codes", backupCodes).Error
}

func (r *UserMFARepository) UpdateVerificationStatus(mfaID uint, verified bool) error {
	db := databases.GetDB()
	return db.Model(&entities.UserMFA{}).Where("id = ?", mfaID).Update("verified", verified).Error
}

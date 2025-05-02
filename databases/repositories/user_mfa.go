package repositories

import (
	"github.com/BlenDMinh/dutgrad-server/databases"
	"github.com/BlenDMinh/dutgrad-server/databases/entities"
)

// UserMFARepository handles database operations for MFA functionality
type UserMFARepository struct {
	*CrudRepository[entities.UserMFA, uint]
}

// NewUserMFARepository creates a new UserMFA repository
func NewUserMFARepository() *UserMFARepository {
	return &UserMFARepository{
		CrudRepository: NewCrudRepository[entities.UserMFA, uint](),
	}
}

// GetByUserID gets MFA data for a specific user
func (r *UserMFARepository) GetByUserID(userID uint) (*entities.UserMFA, error) {
	db := databases.GetDB()
	var mfa entities.UserMFA
	if err := db.Where("user_id = ?", userID).First(&mfa).Error; err != nil {
		return nil, err
	}
	return &mfa, nil
}

// DeleteByUserID deletes MFA data for a specific user
func (r *UserMFARepository) DeleteByUserID(userID uint) error {
	db := databases.GetDB()
	return db.Where("user_id = ?", userID).Delete(&entities.UserMFA{}).Error
}

// UpdateBackupCodes updates the backup codes for a user's MFA
func (r *UserMFARepository) UpdateBackupCodes(mfaID uint, backupCodes entities.BackupCodes) error {
	db := databases.GetDB()
	return db.Model(&entities.UserMFA{}).Where("id = ?", mfaID).Update("backup_codes", backupCodes).Error
}

// UpdateVerificationStatus updates the verification status of MFA for a user
func (r *UserMFARepository) UpdateVerificationStatus(mfaID uint, verified bool) error {
	db := databases.GetDB()
	return db.Model(&entities.UserMFA{}).Where("id = ?", mfaID).Update("verified", verified).Error
}

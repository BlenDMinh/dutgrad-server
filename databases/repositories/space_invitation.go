package repositories

import (
	"github.com/BlenDMinh/dutgrad-server/databases"
	"github.com/BlenDMinh/dutgrad-server/databases/entities"
)

type SpaceInvitationRepository interface {
	ICrudRepository[entities.SpaceInvitation, uint]
	AcceptInvitation(invitationId uint, userId uint) error
	RejectInvitation(invitationId uint, userId uint) error
	CancelInvitation(spaceID uint, invitedUserID uint) error
	CountInvitationByUserID(userID uint) (int64, error)
}

type spaceInvitationRepositoryImpl struct {
	*CrudRepository[entities.SpaceInvitation, uint]
}

func NewSpaceInvitationRepository() SpaceInvitationRepository {
	return &spaceInvitationRepositoryImpl{
		CrudRepository: NewCrudRepository[entities.SpaceInvitation, uint](),
	}
}

func (r *spaceInvitationRepositoryImpl) AcceptInvitation(invitationId uint, userId uint) error {
	var invitation entities.SpaceInvitation
	db := databases.GetDB()
	if err := db.First(&invitation, "id = ? AND invited_user_id = ?", invitationId, userId).Error; err != nil {
		return err
	}

	member := entities.SpaceUser{
		UserID:      userId,
		SpaceID:     invitation.SpaceID,
		SpaceRoleID: &invitation.SpaceRoleID,
	}
	if err := db.Create(&member).Error; err != nil {
		return err
	}

	if err := db.Delete(&invitation).Error; err != nil {
		return err
	}

	return nil
}

func (r *spaceInvitationRepositoryImpl) RejectInvitation(invitationId uint, userId uint) error {
	db := databases.GetDB()

	if err := db.Where("id = ? AND invited_user_id = ?", invitationId, userId).Delete(&entities.SpaceInvitation{}).Error; err != nil {
		return err
	}

	return nil
}

func (r *spaceInvitationRepositoryImpl) CancelInvitation(spaceID uint, invitedUserID uint) error {
	db := databases.GetDB()
	if err := db.Where("space_id = ? AND invited_user_id = ?", spaceID, invitedUserID).
		Unscoped().Delete(&entities.SpaceInvitation{}).Error; err != nil {
		return err
	}
	return nil
}

func (s *spaceInvitationRepositoryImpl) CountInvitationByUserID(userID uint) (int64, error) {
	var count int64
	err := databases.GetDB().
		Model(&entities.SpaceInvitation{}).
		Where("invited_user_id = ? AND status = ?", userID, entities.InvitationStatusPending).
		Count(&count).Error
	return count, err
}

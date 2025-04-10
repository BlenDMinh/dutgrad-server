package repositories

import (
	"github.com/BlenDMinh/dutgrad-server/databases"
	"github.com/BlenDMinh/dutgrad-server/databases/entities"
)

type SpaceRepository struct {
	*CrudRepository[entities.Space, uint]
}

func NewSpaceRepository() *SpaceRepository {
	return &SpaceRepository{
		CrudRepository: NewCrudRepository[entities.Space, uint](),
	}
}

func (r *SpaceRepository) FindPublicSpaces() ([]entities.Space, error) {
	var spaces []entities.Space
	db := databases.GetDB()
	err := db.Where("privacy_status = ?", false).Find(&spaces).Error
	return spaces, err
}

func (r *SpaceRepository) GetMembers(spaceId uint) ([]entities.SpaceUser, error) {
	var members []entities.SpaceUser
	db := databases.GetDB()
	err := db.Preload("User").Preload("SpaceRole").Where("space_id = ?", spaceId).Find(&members).Error
	return members, err
}

func (r *SpaceRepository) GetInvitations(spaceId uint) ([]entities.SpaceInvitation, error) {
	var invitations []entities.SpaceInvitation
	db := databases.GetDB()
	err := db.Preload("InvitedUser").Preload("SpaceRole").Where("space_id = ?", spaceId).Find(&invitations).Error
	return invitations, err
}

func (r *SpaceRepository) GetUserRole(userID, spaceID uint) (*entities.SpaceRole, error) {
	db := databases.GetDB()
	var spaceUser entities.SpaceUser

	// Find the space user relation with preloaded role
	err := db.Where("user_id = ? AND space_id = ?", userID, spaceID).
		Preload("SpaceRole").
		First(&spaceUser).Error

	if err != nil {
		return nil, err
	}

	if spaceUser.SpaceRoleID == nil {
		return nil, nil
	}

	return &spaceUser.SpaceRole, nil
}

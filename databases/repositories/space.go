package repositories

import (
	"errors"

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

func (s *SpaceRepository) CreateInvitation(invitation *entities.SpaceInvitation) (*entities.SpaceInvitation, error) {
	db := databases.GetDB()
	if err := db.Create(invitation).Error; err != nil {
		return nil, err
	}
	return invitation, nil
}

func (r *SpaceRepository) GetAllRoles() ([]entities.SpaceRole, error) {
	var roles []entities.SpaceRole
	db := databases.GetDB()
	if err := db.Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

func (s *SpaceRepository) JoinPublicSpace(spaceID uint, userID uint) error {
	db := databases.GetDB()
	var space entities.Space
	if err := db.First(&space, spaceID).Error; err != nil {
		return errors.New("space not found")
	}
	var count int64
	db.Model(&entities.SpaceUser{}).
		Where("space_id = ? AND user_id = ?", spaceID, userID).
		Count(&count)
	if count > 0 {
		return errors.New("user is already a member of this space")
	}

	roleID := uint(entities.Viewer)
	spaceUser := entities.SpaceUser{
		UserID:      userID,
		SpaceID:     spaceID,
		SpaceRoleID: &roleID,
	}
	if err := db.Create(&spaceUser).Error; err != nil {
		return err
	}

	return nil
}

func (s *SpaceRepository) IsMemberOfSpace(userID uint, spaceID uint) (bool, error) {
	var count int64
	db := databases.GetDB()
	err := db.Model(&entities.SpaceUser{}).
		Where("user_id = ? AND space_id = ?", userID, spaceID).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (s *SpaceRepository) CountSpacesByUserID(userID uint) (int64, error) {
	var count int64
	err := databases.GetDB().
		Table("space_users").
		Where("space_users.user_id = ?", userID).
		Count(&count).Error
	return count, err
}

func (r *SpaceRepository) GetPopularSpaces(order string) ([]entities.Space, error) {
	var spaces []entities.Space
	db := databases.GetDB()

	if order == "member_count" {
		err := db.Model(&entities.Space{}).
			Select("spaces.*, COUNT(space_users.user_id) as member_count").
			Where("privacy_status = ?", false).
			Joins("LEFT JOIN space_users ON space_users.space_id = spaces.id").
			Group("spaces.id").
			Order("member_count DESC").
			Find(&spaces).Error
		return spaces, err
	}

	return []entities.Space{}, nil
}

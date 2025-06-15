package repositories

import (
	"errors"
	"time"

	"github.com/BlenDMinh/dutgrad-server/databases"
	"github.com/BlenDMinh/dutgrad-server/databases/entities"
	"github.com/BlenDMinh/dutgrad-server/models/dtos"
)

type SpaceRepository interface {
	ICrudRepository[entities.Space, uint]
	FindPublicSpaces(page int, pageSize int) ([]*entities.Space, error)
	CountPublicSpaces() (int64, error)
	GetMembers(spaceId uint) ([]entities.SpaceUser, error)
	GetInvitations(spaceId uint) ([]entities.SpaceInvitation, error)
	GetUserRole(userID, spaceID uint) (*entities.SpaceRole, error)
	CreateInvitation(invitation *entities.SpaceInvitation) (*entities.SpaceInvitation, error)
	GetAllRoles() ([]entities.SpaceRole, error)
	JoinPublicSpace(spaceID uint, userID uint) error
	IsMemberOfSpace(userID uint, spaceID uint) (bool, error)
	CountSpacesByUserID(userID uint) (int64, error)
	CountOwnedSpacesByUserID(userID uint) (int64, error)
	GetPopularSpaces(order string) ([]*entities.Space, error)
	UpdateMemberRole(spaceID, memberID, roleID, updatedBy uint) error
	RemoveMember(spaceID, memberID uint) error
	GetSpaceUsage(spaceID uint) (*dtos.SpaceUsage, error)
}

type spaceRepositoryImpl struct {
	*CrudRepository[entities.Space, uint]
}

func NewSpaceRepository() SpaceRepository {
	return &spaceRepositoryImpl{
		CrudRepository: NewCrudRepository[entities.Space, uint](),
	}
}

func (r *spaceRepositoryImpl) FindPublicSpaces(page int, pageSize int) ([]*entities.Space, error) {
	var spaces []*entities.Space
	db := databases.GetDB()

	pagination := NewPagination(page, pageSize, DefaultPageSize)

	err := pagination.ApplyPagination(db).Where("privacy_status = ?", false).Find(&spaces).Error
	if err != nil {
		return nil, err
	}

	spaces, err = r.aggregateUserCount(spaces)
	if err != nil {
		return nil, err
	}

	return spaces, err
}

func (r *spaceRepositoryImpl) CountPublicSpaces() (int64, error) {
	var count int64
	db := databases.GetDB()
	err := db.Model(&entities.Space{}).Where("privacy_status = ?", false).Count(&count).Error
	return count, err
}

func (r *spaceRepositoryImpl) GetMembers(spaceId uint) ([]entities.SpaceUser, error) {
	var members []entities.SpaceUser
	db := databases.GetDB()
	err := db.Preload("User").Preload("SpaceRole").Where("space_id = ?", spaceId).Find(&members).Error
	return members, err
}

func (r *spaceRepositoryImpl) GetInvitations(spaceId uint) ([]entities.SpaceInvitation, error) {
	var invitations []entities.SpaceInvitation
	db := databases.GetDB()
	err := db.Preload("InvitedUser").Preload("SpaceRole").Where("space_id = ?", spaceId).Find(&invitations).Error
	return invitations, err
}

func (r *spaceRepositoryImpl) GetUserRole(userID, spaceID uint) (*entities.SpaceRole, error) {
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

func (s *spaceRepositoryImpl) CreateInvitation(invitation *entities.SpaceInvitation) (*entities.SpaceInvitation, error) {
	db := databases.GetDB()
	if err := db.Create(invitation).Error; err != nil {
		return nil, err
	}
	return invitation, nil
}

func (r *spaceRepositoryImpl) GetAllRoles() ([]entities.SpaceRole, error) {
	var roles []entities.SpaceRole
	db := databases.GetDB()
	if err := db.Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

func (s *spaceRepositoryImpl) JoinPublicSpace(spaceID uint, userID uint) error {
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

	roleID := uint(entities.SpaceRoleViewer)
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

func (s *spaceRepositoryImpl) IsMemberOfSpace(userID uint, spaceID uint) (bool, error) {
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

func (s *spaceRepositoryImpl) CountSpacesByUserID(userID uint) (int64, error) {
	var count int64
	err := databases.GetDB().
		Table("space_users").
		Where("space_users.user_id = ?", userID).
		Count(&count).Error
	return count, err
}

func (s *spaceRepositoryImpl) CountOwnedSpacesByUserID(userID uint) (int64, error) {
	var count int64
	err := databases.GetDB().
		Table("space_users").
		Where("space_users.user_id = ? AND space_users.space_role_id = ?", userID, entities.SpaceRoleOwner).
		Count(&count).Error
	return count, err
}

func (r *spaceRepositoryImpl) GetPopularSpaces(order string) ([]*entities.Space, error) {
	var spaces []*entities.Space
	db := databases.GetDB()

	if order == "user_count" {
		err := db.Model(&entities.Space{}).
			Select("spaces.*, COUNT(space_users.user_id) as user_count").
			Where("privacy_status = ?", false).
			Joins("LEFT JOIN space_users ON space_users.space_id = spaces.id").
			Group("spaces.id").
			Order("user_count DESC").
			Find(&spaces).Error
		spaces, err := r.aggregateUserCount(spaces)
		if err != nil {
			return nil, err
		}

		return spaces, err
	}

	spaces, err := r.aggregateUserCount(spaces)
	if err != nil {
		return nil, err
	}

	return []*entities.Space{}, nil
}

func (r *spaceRepositoryImpl) UpdateMemberRole(spaceID, memberID, roleID, updatedBy uint) error {
	db := databases.GetDB()

	var spaceUser entities.SpaceUser
	if err := db.Where("space_id = ? AND user_id = ?", spaceID, memberID).First(&spaceUser).Error; err != nil {
		return errors.New("member not found in the space")
	}

	spaceUser.SpaceRoleID = &roleID
	if err := db.Save(&spaceUser).Error; err != nil {
		return errors.New("failed to update member role")
	}

	return nil
}

func (r *spaceRepositoryImpl) RemoveMember(spaceID, memberID uint) error {
	db := databases.GetDB()

	var spaceUser entities.SpaceUser
	if err := db.Where("space_id = ? AND user_id = ?", spaceID, memberID).First(&spaceUser).Error; err != nil {
		return errors.New("member not found in the space")
	}

	if err := db.Delete(&spaceUser).Error; err != nil {
		return errors.New("failed to remove member from the space")
	}

	return nil
}

func (r *spaceRepositoryImpl) GetSpaceUsage(spaceID uint) (*dtos.SpaceUsage, error) {
	db := databases.GetDB()
	var usage dtos.SpaceUsage
	usage.SpaceID = spaceID

	today := time.Now().Format("2006-01-02")

	err := db.Model(&entities.ChatHistory{}).
		Joins("JOIN user_query_sessions ON user_query_sessions.id = chat_histories.session_id").
		Where("user_query_sessions.space_id = ? AND DATE(chat_histories.created_at) = ?", spaceID, today).
		Count(&usage.ChatAPICallsUsageDaily).Error

	if err != nil {
		return nil, err
	}

	return &usage, nil
}

func (r *spaceRepositoryImpl) aggregateUserCount(spaces []*entities.Space) ([]*entities.Space, error) {
	type SpaceWithUserCount struct {
		SpaceID   uint
		UserCount int64
	}

	spaceIds := make([]uint, len(spaces))
	for i, space := range spaces {
		spaceIds[i] = space.ID
	}

	var spaceCounts []SpaceWithUserCount
	db := databases.GetDB()
	err := db.Table("space_users").
		Select("space_id, COUNT(user_id) as user_count").
		Where("space_id IN (?)", spaceIds).
		Group("space_id").
		Scan(&spaceCounts).Error
	if err != nil {
		return nil, err
	}

	spaceCountsMap := make(map[uint]int64)
	for _, sc := range spaceCounts {
		spaceCountsMap[sc.SpaceID] = sc.UserCount
	}

	for _, space := range spaces {
		if count, exists := spaceCountsMap[space.ID]; exists {
			space.UserCount = int(count)
		} else {
			space.UserCount = 0
		}
	}

	return spaces, nil
}

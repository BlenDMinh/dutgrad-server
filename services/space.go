package services

import (
	"fmt"

	"github.com/BlenDMinh/dutgrad-server/databases"
	"github.com/BlenDMinh/dutgrad-server/databases/entities"
	"github.com/BlenDMinh/dutgrad-server/databases/repositories"
	"github.com/BlenDMinh/dutgrad-server/helpers"
)

type SpaceService struct {
	CrudService[entities.Space, uint]
	invitationLinkRepo repositories.SpaceInvitationLinkRepository
}

func NewSpaceService(invitationLinkRepo repositories.SpaceInvitationLinkRepository) *SpaceService {
	return &SpaceService{
		CrudService:        *NewCrudService(repositories.NewSpaceRepository()),
		invitationLinkRepo: invitationLinkRepo,
	}
}

func (s *SpaceService) GetPublicSpaces() ([]entities.Space, error) {
	return s.repo.(*repositories.SpaceRepository).FindPublicSpaces()
}

func (s *SpaceService) GetMembers(spaceId uint) ([]entities.SpaceUser, error) {
	return s.repo.(*repositories.SpaceRepository).GetMembers(spaceId)
}

func (s *SpaceService) GetInvitations(spaceId uint) ([]entities.SpaceInvitation, error) {
	return s.repo.(*repositories.SpaceRepository).GetInvitations(spaceId)
}

func (s *SpaceService) GetOrCreateSpaceInvitationLink(spaceID, spaceRoleID uint) (*entities.SpaceInvitationLink, error) {
	repo := s.invitationLinkRepo
	invitationLink, _ := repo.GetBySpaceID(spaceID)
	if invitationLink == nil {
		invitationLink = &entities.SpaceInvitationLink{
			SpaceID:     spaceID,
			SpaceRoleID: spaceRoleID,
		}
		invitationLink, err := repo.Create(invitationLink)

		if err != nil {
			return nil, err
		}
		return invitationLink, nil
	}
	if invitationLink.SpaceRoleID != spaceRoleID {
		invitationLink.SpaceRoleID = spaceRoleID
		invitationLink, err := repo.Update(invitationLink)
		if err != nil {
			return nil, err
		}

		return invitationLink, nil
	}
	return invitationLink, nil
}

func (s *SpaceService) CreateInvitation(invitation *entities.SpaceInvitation) (*entities.SpaceInvitation, error) {
	return s.repo.(*repositories.SpaceRepository).CreateInvitation(invitation)
}

func (s *SpaceService) GetSpaceRoles() ([]entities.SpaceRole, error) {
	return s.repo.(*repositories.SpaceRepository).GetAllRoles()
}

func (s *SpaceService) JoinSpaceWithToken(token string, userID uint) (uint, error) {
	payload, err := helpers.VerifyTokenForPayload(token)
	if err != nil {
		return 0, err
	}

	if payload == nil {
		return 0, fmt.Errorf("invalid token")
	}

	parsePayload := *payload

	db := databases.GetDB()
	var spaceUser entities.SpaceUser
	err = db.Where("user_id = ? AND space_id = ?", userID, parsePayload["space_id"]).First(&spaceUser).Error
	if err == nil {
		return 0, fmt.Errorf("user is already a member of this space")
	}

	spaceRoleIDFloat := parsePayload["space_role_id"].(float64)
	spaceRoleID := uint(spaceRoleIDFloat)

	spaceIDFloat := parsePayload["space_id"].(float64)
	spaceID := uint(spaceIDFloat)

	newSpaceUser := entities.SpaceUser{
		UserID:      userID,
		SpaceID:     spaceID,
		SpaceRoleID: &spaceRoleID,
	}

	err = db.Create(&newSpaceUser).Error
	if err != nil {
		return 0, err
	}

	return spaceID, nil
}

func (s *SpaceService) GetUserRole(userID, spaceID uint) (*entities.SpaceRole, error) {
	role, err := s.repo.(*repositories.SpaceRepository).GetUserRole(userID, spaceID)
	if err != nil {
		return nil, fmt.Errorf("user is not a member of this space or %v", err)
	}
	if role == nil {
		return nil, fmt.Errorf("user has no role in this space")
	}
	return role, nil
}

func (s *SpaceService) JoinPublicSpace(spaceID uint, userID uint) error {
	return s.repo.(*repositories.SpaceRepository).JoinPublicSpace(spaceID, userID)
}

func (s *SpaceService) IsMemberOfSpace(userID uint, spaceID uint) (bool, error) {
	return s.repo.(*repositories.SpaceRepository).IsMemberOfSpace(userID, spaceID)
}

func (s *SpaceService) CountSpacesByUserID(userID uint) (int64, error) {
	return s.repo.(*repositories.SpaceRepository).CountSpacesByUserID(userID)

}

func (s *SpaceService) GetPopularSpaces(order string) ([]entities.Space, error) {
	return s.repo.(*repositories.SpaceRepository).GetPopularSpaces(order)
}

func (s *SpaceService) CheckSpaceCreationLimit(userID uint) error {
	db := databases.GetDB()

	var user entities.User
	if err := db.Preload("Tier").Where("id = ?", userID).First(&user).Error; err != nil {
		return err
	}

	count, err := s.CountSpacesByUserID(userID)
	if err != nil {
		return err
	}

	spaceLimit := 5

	if user.Tier != nil {
		spaceLimit = user.Tier.SpaceLimit
	}

	if count >= int64(spaceLimit) {
		return fmt.Errorf("space limit reached: you can only create %d spaces with your current tier", spaceLimit)
	}

	return nil
}

func (s *SpaceService) CreateSpace(space *entities.Space, userID uint) (*entities.Space, error) {
	if err := s.CheckSpaceCreationLimit(userID); err != nil {
		return nil, err
	}

	createdSpace, err := s.Create(space)
	if err != nil {
		return nil, err
	}

	ownerRoleID := uint(entities.Owner)

	spaceUser := entities.SpaceUser{
		UserID:      userID,
		SpaceID:     createdSpace.ID,
		SpaceRoleID: &ownerRoleID,
	}

	db := databases.GetDB()
	if err := db.Create(&spaceUser).Error; err != nil {
		return nil, fmt.Errorf("failed to add user as owner: %v", err)
	}

	return createdSpace, nil
}

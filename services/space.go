package services

import (
	"errors"
	"fmt"

	"github.com/BlenDMinh/dutgrad-server/databases"
	"github.com/BlenDMinh/dutgrad-server/databases/entities"
	"github.com/BlenDMinh/dutgrad-server/databases/repositories"
	"github.com/BlenDMinh/dutgrad-server/helpers"
	"github.com/BlenDMinh/dutgrad-server/models/dtos"
)

type SpaceService interface {
	ICrudService[entities.Space, uint]
	GetPublicSpaces(page int, pageSize int) (*helpers.PaginationResult, error)
	GetMembers(spaceId uint) ([]entities.SpaceUser, error)
	GetInvitations(spaceId uint) ([]entities.SpaceInvitation, error)
	GetOrCreateSpaceInvitationLink(spaceID, spaceRoleID uint) (*entities.SpaceInvitationLink, error)
	CreateInvitation(invitation *entities.SpaceInvitation) (*entities.SpaceInvitation, error)
	GetSpaceRoles() ([]entities.SpaceRole, error)
	JoinSpaceWithToken(token string, userID uint) (uint, error)
	JoinPublicSpace(spaceID uint, userID uint) error
	GetUserRole(userID, spaceID uint) (*entities.SpaceRole, error)
	IsMemberOfSpace(userID uint, spaceID uint) (bool, error)
	CountSpacesByUserID(userID uint) (int64, error)
	CountSpaceMembers(spaceID uint) (int64, error)
	GetPopularSpaces(order string) ([]entities.Space, error)
	CheckSpaceCreationLimit(userID uint) error
	CreateSpace(space *entities.Space, userID uint) (*entities.Space, error)
	UpdateMemberRole(spaceID, memberID, roleID, updatedBy uint) error
	RemoveMember(spaceID, memberID, requestingUserID uint) error
	Delete(id uint) error
	GetSpaceUsage(spaceID uint) (*dtos.SpaceUsage, error)
	IsAPIRateLimited(spaceID uint) bool
}

type spaceServiceImpl struct {
	CrudService[entities.Space, uint]
	repo                      repositories.SpaceRepository
	invitationLinkRepo        repositories.SpaceInvitationLinkRepository
	ragServerService          *RAGServerService
	userRepository            repositories.UserRepository
	spaceInvitationRepository repositories.SpaceInvitationRepository
	documentRepository        repositories.DocumentRepository
}

func NewSpaceService(
	invitationLinkRepo repositories.SpaceInvitationLinkRepository,
	ragServerService *RAGServerService,
	userRepository repositories.UserRepository,
	spaceInvitationRepository repositories.SpaceInvitationRepository,
	documentRepository repositories.DocumentRepository,
) SpaceService {
	crudService := NewCrudService(repositories.NewSpaceRepository())
	repo := crudService.repo.(repositories.SpaceRepository)

	return &spaceServiceImpl{
		CrudService:               *crudService,
		invitationLinkRepo:        invitationLinkRepo,
		ragServerService:          ragServerService,
		repo:                      repo,
		userRepository:            userRepository,
		spaceInvitationRepository: spaceInvitationRepository,
		documentRepository:        documentRepository,
	}
}

func (s *spaceServiceImpl) GetPublicSpaces(page int, pageSize int) (*helpers.PaginationResult, error) {
	spaces, err := s.repo.FindPublicSpaces(page, pageSize)
	if err != nil {
		return nil, err
	}

	count, err := s.repo.CountPublicSpaces()
	if err != nil {
		return nil, err
	}

	result := helpers.CreatePaginationResult(spaces, page, pageSize, count)
	return &result, nil
}

func (s *spaceServiceImpl) GetMembers(spaceId uint) ([]entities.SpaceUser, error) {
	return s.repo.GetMembers(spaceId)
}

func (s *spaceServiceImpl) GetInvitations(spaceId uint) ([]entities.SpaceInvitation, error) {
	return s.repo.GetInvitations(spaceId)
}

func (s *spaceServiceImpl) GetOrCreateSpaceInvitationLink(spaceID, spaceRoleID uint) (*entities.SpaceInvitationLink, error) {
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

func (s *spaceServiceImpl) CreateInvitation(invitation *entities.SpaceInvitation) (*entities.SpaceInvitation, error) {
	return s.repo.CreateInvitation(invitation)
}

func (s *spaceServiceImpl) GetSpaceRoles() ([]entities.SpaceRole, error) {
	return s.repo.GetAllRoles()
}

func (s *spaceServiceImpl) JoinSpaceWithToken(token string, userID uint) (uint, error) {
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

func (s *spaceServiceImpl) GetUserRole(userID, spaceID uint) (*entities.SpaceRole, error) {
	role, err := s.repo.GetUserRole(userID, spaceID)
	if err != nil {
		return nil, fmt.Errorf("user is not a member of this space or %v", err)
	}
	if role == nil {
		return nil, fmt.Errorf("user has no role in this space")
	}
	return role, nil
}

func (s *spaceServiceImpl) JoinPublicSpace(spaceID uint, userID uint) error {
	return s.repo.JoinPublicSpace(spaceID, userID)
}

func (s *spaceServiceImpl) IsMemberOfSpace(userID uint, spaceID uint) (bool, error) {
	return s.repo.IsMemberOfSpace(userID, spaceID)
}

func (s *spaceServiceImpl) CountSpacesByUserID(userID uint) (int64, error) {
	return s.repo.CountSpacesByUserID(userID)
}

func (s *spaceServiceImpl) CountSpaceMembers(spaceID uint) (int64, error) {
	members, err := s.GetMembers(spaceID)
	if err != nil {
		return 0, err
	}
	return int64(len(members)), nil
}

func (s *spaceServiceImpl) GetPopularSpaces(order string) ([]entities.Space, error) {
	return s.repo.GetPopularSpaces(order)
}

func (s *spaceServiceImpl) CheckSpaceCreationLimit(userID uint) error {
	user, err := s.userRepository.GetById(userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %v", err)
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

func (s *spaceServiceImpl) CreateSpace(space *entities.Space, userID uint) (*entities.Space, error) {
	if err := s.CheckSpaceCreationLimit(userID); err != nil {
		return nil, err
	}

	createdSpace, err := s.Create(space)
	if err != nil {
		return nil, err
	}

	ownerRoleID := uint(entities.SpaceRoleOwner)

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

func (s *spaceServiceImpl) UpdateMemberRole(spaceID, memberID, roleID, updatedBy uint) error {
	return s.repo.UpdateMemberRole(spaceID, memberID, roleID, updatedBy)
}

func (s *spaceServiceImpl) RemoveMember(spaceID, memberID, requestingUserID uint) error {
	requestingUserRole, err := s.GetUserRole(requestingUserID, spaceID)
	if err != nil {
		return err
	}

	if !requestingUserRole.IsOwner() {
		return errors.New("only space owners can remove members")
	}

	if memberID == requestingUserID {
		return errors.New("you cannot remove yourself from the space")
	}

	isMember, err := s.repo.IsMemberOfSpace(memberID, spaceID)
	if err != nil {
		return err
	}

	if isMember {
		memberRole, err := s.GetUserRole(memberID, spaceID)
		if err != nil {
			return err
		}

		if memberRole.ID == uint(entities.SpaceRoleOwner) {
			return errors.New("cannot remove a space owner")
		}

		return s.repo.RemoveMember(spaceID, memberID)
	}

	return s.spaceInvitationRepository.CancelInvitation(spaceID, memberID)
}

func (s *spaceServiceImpl) Delete(id uint) error {
	documents, err := s.documentRepository.GetBySpaceID(id)
	if err != nil {
		return fmt.Errorf("failed to get documents in space: %v", err)
	}

	for _, doc := range documents {
		err := s.documentRepository.Delete(doc.ID)
		if err != nil {
			return fmt.Errorf("failed to delete document %d: %v", doc.ID, err)
		}
	}

	err = s.ragServerService.RemoveSpace(id)
	if err != nil {
		return fmt.Errorf("failed to remove space from RAG server: %v", err)
	}

	return s.repo.Delete(id)
}

func (s *spaceServiceImpl) GetSpaceUsage(spaceID uint) (*dtos.SpaceUsage, error) {
	usage, err := s.repo.GetSpaceUsage(spaceID)
	if err != nil {
		return nil, fmt.Errorf("failed to get space usage: %v", err)
	}

	if usage == nil {
		return &dtos.SpaceUsage{
			SpaceID:                spaceID,
			ChatAPICallsUsageDaily: 0,
		}, nil
	}

	return usage, nil
}

func (s *spaceServiceImpl) IsAPIRateLimited(spaceID uint) bool {
	space, err := s.repo.GetById(spaceID)
	if err != nil {
		return false
	}
	usage, err := s.GetSpaceUsage(spaceID)
	if err != nil {
		return false
	}

	return usage.ChatAPICallsUsageDaily >= int64(space.ApiCallLimit)
}

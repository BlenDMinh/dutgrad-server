package services

import (
	"github.com/BlenDMinh/dutgrad-server/databases/entities"
	"github.com/BlenDMinh/dutgrad-server/databases/repositories"
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
	// Update the SpaceRoleID if it is different
	if invitationLink.SpaceRoleID != spaceRoleID {
		invitationLink.SpaceRoleID = spaceRoleID
		invitationLink, err := repo.Update(invitationLink)
		if err != nil {
			return nil, err
		}

		return invitationLink, nil
	}
	// If the invitation link already exists, return it
	return invitationLink, nil
}

func (s *SpaceService) CreateInvitation(invitation *entities.SpaceInvitation) (*entities.SpaceInvitation, error) {
	return s.repo.(*repositories.SpaceRepository).CreateInvitation(invitation)
}

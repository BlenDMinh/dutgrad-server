package services

import (
	"github.com/BlenDMinh/dutgrad-server/databases/entities"
	"github.com/BlenDMinh/dutgrad-server/databases/repositories"
)

type SpaceInvitationService struct {
	CrudService[entities.SpaceInvitation, uint]
}

func NewSpaceInvitationService() *SpaceInvitationService {
	return &SpaceInvitationService{
		CrudService: *NewCrudService(repositories.NewSpaceInvitationRepository()),
	}
}

func (s *SpaceInvitationService) AcceptInvitation(invitationId uint, userId uint) error {
	return s.repo.(*repositories.SpaceInvitationRepository).AcceptInvitation(invitationId, userId)
}

func (s *SpaceInvitationService) RejectInvitation(invitationId uint, userId uint) error {
	return s.repo.(*repositories.SpaceInvitationRepository).RejectInvitation(invitationId, userId)
}

func (s *SpaceInvitationService) CancelInvitation(spaceID uint, invitedUserID uint) error {
	return s.repo.(*repositories.SpaceInvitationRepository).CancelInvitation(spaceID, invitedUserID)
}

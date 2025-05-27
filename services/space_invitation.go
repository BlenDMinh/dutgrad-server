package services

import (
	"github.com/BlenDMinh/dutgrad-server/databases/entities"
	"github.com/BlenDMinh/dutgrad-server/databases/repositories"
)

type SpaceInvitationService interface {
	ICrudService[entities.SpaceInvitation, uint]
	AcceptInvitation(invitationId uint, userId uint) error
	RejectInvitation(invitationId uint, userId uint) error
	CancelInvitation(spaceID uint, invitedUserID uint) error
	CountInvitationByUserID(userID uint) (int64, error)
}

type spaceInvitationServiceImpl struct {
	CrudService[entities.SpaceInvitation, uint]
	repo repositories.SpaceInvitationRepository
}

func NewSpaceInvitationService() SpaceInvitationService {
	crudService := NewCrudService(repositories.NewSpaceInvitationRepository())
	repo := crudService.repo.(repositories.SpaceInvitationRepository)
	return &spaceInvitationServiceImpl{
		CrudService: *crudService,
		repo:        repo,
	}
}

func (s *spaceInvitationServiceImpl) AcceptInvitation(invitationId uint, userId uint) error {
	return s.repo.AcceptInvitation(invitationId, userId)
}

func (s *spaceInvitationServiceImpl) RejectInvitation(invitationId uint, userId uint) error {
	return s.repo.RejectInvitation(invitationId, userId)
}

func (s *spaceInvitationServiceImpl) CancelInvitation(spaceID uint, invitedUserID uint) error {
	return s.repo.CancelInvitation(spaceID, invitedUserID)
}

func (s *spaceInvitationServiceImpl) CountInvitationByUserID(userID uint) (int64, error) {
	return s.repo.CountInvitationByUserID(userID)
}

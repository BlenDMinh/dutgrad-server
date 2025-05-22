package repositories

import (
	"github.com/BlenDMinh/dutgrad-server/databases"
	"github.com/BlenDMinh/dutgrad-server/databases/entities"
)

type SpaceInvitationLinkRepository interface {
	ICrudRepository[entities.SpaceInvitationLink, uint]
	GetBySpaceID(spaceID uint) (*entities.SpaceInvitationLink, error)
}

type spaceInvitationLinkRepositoryImpl struct {
	*CrudRepository[entities.SpaceInvitationLink, uint]
}

func NewSpaceInvitationLinkRepository() SpaceInvitationLinkRepository {
	return &spaceInvitationLinkRepositoryImpl{
		CrudRepository: NewCrudRepository[entities.SpaceInvitationLink, uint](),
	}
}

func (s *spaceInvitationLinkRepositoryImpl) GetBySpaceID(spaceID uint) (*entities.SpaceInvitationLink, error) {
	var invitationLink entities.SpaceInvitationLink
	db := databases.GetDB()
	if err := db.Where("space_id = ?", spaceID).First(&invitationLink).Error; err != nil {
		return nil, err
	}
	return &invitationLink, nil
}

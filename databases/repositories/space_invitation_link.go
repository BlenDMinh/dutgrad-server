package repositories

import (
	"github.com/BlenDMinh/dutgrad-server/databases"
	"github.com/BlenDMinh/dutgrad-server/databases/entities"
)

type SpaceInvitationLinkRepository struct {
	*CrudRepository[entities.SpaceInvitationLink, uint]
}

func NewSpaceInvitationLinkRepository() *SpaceInvitationLinkRepository {
	return &SpaceInvitationLinkRepository{
		CrudRepository: NewCrudRepository[entities.SpaceInvitationLink, uint](),
	}
}

func (s *SpaceInvitationLinkRepository) GetBySpaceID(spaceID uint) (*entities.SpaceInvitationLink, error) {
	var invitationLink entities.SpaceInvitationLink
	db := databases.GetDB()
	if err := db.Where("space_id = ?", spaceID).First(&invitationLink).Error; err != nil {
		return nil, err
	}
	return &invitationLink, nil
}

package repositories

import "github.com/BlenDMinh/dutgrad-server/databases/entities"

type SpaceInvitationLinkRepository struct {
	*CrudRepository[entities.SpaceInvitationLink, uint]
}

func NewSpaceInvitationLinkRepository() *SpaceInvitationLinkRepository {
	return &SpaceInvitationLinkRepository{
		CrudRepository: NewCrudRepository[entities.SpaceInvitationLink, uint](),
	}
}

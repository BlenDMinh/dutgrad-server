package repositories

import "github.com/BlenDMinh/dutgrad-server/databases/entities"

type SpaceInvitationRepository struct {
	*CrudRepository[entities.SpaceInvitation, uint]
}

func NewSpaceInvitationRepository() *SpaceInvitationRepository {
	return &SpaceInvitationRepository{
		CrudRepository: NewCrudRepository[entities.SpaceInvitation, uint](),
	}
}

package services

import (
	"github.com/BlenDMinh/dutgrad-server/databases"
	"github.com/BlenDMinh/dutgrad-server/databases/entities"
	"github.com/BlenDMinh/dutgrad-server/databases/repositories"
	"github.com/BlenDMinh/dutgrad-server/models/dtos"
)

type UserService struct {
	CrudService[entities.User, uint]
}

func NewUserService() *UserService {
	return &UserService{
		CrudService: *NewCrudService(repositories.NewUserRepository()),
	}
}

func (s *UserService) GetSpacesByUserId(userId uint) ([]entities.Space, error) {
	return s.repo.(*repositories.UserRepository).GetSpacesByUserId(userId)
}

func (s *UserService) GetUserByEmail(email string) (*entities.User, error) {
	return s.repo.(*repositories.UserRepository).GetUserByEmail(email)
}

func (s *UserService) GetInvitationsByUserId(InvitedUserId uint) ([]entities.SpaceInvitation, error) {
	return s.repo.(*repositories.UserRepository).GetInvitationsByUserId(InvitedUserId)
}

func (s *UserService) SearchUsers(query string) ([]entities.User, error) {
	return s.repo.(*repositories.UserRepository).SearchUsers(query)
}

func (s *UserService) GetUserTier(userID uint) (*entities.Tier, error) {
	db := databases.GetDB()
	var user entities.User

	if err := db.Preload("Tier").First(&user, userID).Error; err != nil {
		return nil, err
	}

	return user.Tier, nil
}

func (s *UserService) GetUserTierUsage(userID uint) (*dtos.TierUsageResponse, error) {
	return s.repo.(*repositories.UserRepository).GetUserTierUsage(userID)
}

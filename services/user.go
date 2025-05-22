package services

import (
	"github.com/BlenDMinh/dutgrad-server/databases"
	"github.com/BlenDMinh/dutgrad-server/databases/entities"
	"github.com/BlenDMinh/dutgrad-server/databases/repositories"
	"github.com/BlenDMinh/dutgrad-server/models/dtos"
)

type UserService interface {
	ICrudService[entities.User, uint]
	GetSpacesByUserId(userId uint) ([]entities.Space, error)
	GetUserByEmail(email string) (*entities.User, error)
	GetInvitationsByUserId(InvitedUserId uint) ([]entities.SpaceInvitation, error)
	SearchUsers(query string) ([]entities.User, error)
	GetUserTier(userID uint) (*entities.Tier, error)
	GetUserTierUsage(userID uint) (*dtos.TierUsageResponse, error)
}

type UserServiceImpl struct {
	CrudService[entities.User, uint]
	repo repositories.UserRepository
}

func NewUserService() UserService {
	crudService := NewCrudService(repositories.NewUserRepository())
	repo := crudService.repo.(repositories.UserRepository)

	return &UserServiceImpl{
		CrudService: *crudService,
		repo:        repo,
	}
}

func (s *UserServiceImpl) GetSpacesByUserId(userId uint) ([]entities.Space, error) {
	return s.repo.GetSpacesByUserId(userId)
}

func (s *UserServiceImpl) GetUserByEmail(email string) (*entities.User, error) {
	return s.repo.GetByEmail(email)
}

func (s *UserServiceImpl) GetInvitationsByUserId(InvitedUserId uint) ([]entities.SpaceInvitation, error) {
	return s.repo.GetInvitationsByUserId(InvitedUserId)
}

func (s *UserServiceImpl) SearchUsers(query string) ([]entities.User, error) {
	return s.repo.SearchUsers(query)
}

func (s *UserServiceImpl) GetUserTier(userID uint) (*entities.Tier, error) {
	db := databases.GetDB()
	var user entities.User

	if err := db.Preload("Tier").First(&user, userID).Error; err != nil {
		return nil, err
	}

	return user.Tier, nil
}

func (s *UserServiceImpl) GetUserTierUsage(userID uint) (*dtos.TierUsageResponse, error) {
	return s.repo.GetUserTierUsage(userID)
}

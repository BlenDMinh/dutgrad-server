package repositories

import (
	"strings"

	"github.com/BlenDMinh/dutgrad-server/databases"
	"github.com/BlenDMinh/dutgrad-server/databases/entities"
)

type UserRepository struct {
	*CrudRepository[entities.User, uint]
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		CrudRepository: NewCrudRepository[entities.User, uint](),
	}
}

func (r *UserRepository) GetSpacesByUserId(userId uint) ([]entities.Space, error) {
	var spaces []entities.Space
	db := databases.GetDB()
	err := db.Joins("JOIN space_users ON space_users.space_id = spaces.id").
		Where("space_users.user_id = ?", userId).
		Find(&spaces).Error
	return spaces, err
}

func (r *UserRepository) GetUserByEmail(email string) (*entities.User, error) {
	var user entities.User
	db := databases.GetDB()
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetInvitationsByUserId(InvitedUserId uint) ([]entities.SpaceInvitation, error) {
	var invitations []entities.SpaceInvitation
	db := databases.GetDB()
	err := db.Preload("Space").Preload("Inviter").
		Where("invited_user_id = ?", InvitedUserId).Find(&invitations).Error
	if err != nil {
		return nil, err
	}
	return invitations, nil
}

func (r *UserRepository) SearchUsers(query string) ([]entities.User, error) {
	var users []entities.User
	db := databases.GetDB()

	if strings.Contains(query, "@") {
		err := db.Where("email = ?", query).Find(&users).Error
		if err != nil {
			return nil, err
		}
	} else {
		err := db.Where("username ILIKE ?", "%"+query+"%").Find(&users).Error
		if err != nil {
			return nil, err
		}
	}

	return users, nil
}

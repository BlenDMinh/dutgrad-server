package repositories

import (
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

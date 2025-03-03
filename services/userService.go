package services

import (
	"github.com/BlenDMinh/dutgrad-server/databases"
	"github.com/BlenDMinh/dutgrad-server/databases/entities"
)

type UserService struct{}

func (s *UserService) GetUserByID(userID uint) (*entities.User, error) {
	db := databases.GetDB()
	var user entities.User
	if err := db.First(&user, userID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

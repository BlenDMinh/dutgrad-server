package seeders

import (
	"github.com/BlenDMinh/dutgrad-server/databases"
	"github.com/BlenDMinh/dutgrad-server/databases/entities"
	"gorm.io/gorm"
)

type SpaceRoleSeeder struct{}

func (s *SpaceRoleSeeder) Name() string {
	return "SpaceRoleSeeder"
}

func (s *SpaceRoleSeeder) Seed() error {
	dbtx := databases.GetDB()

	spaceRoles := []entities.SpaceRole{
		{
			ID:         1,
			Name:       "owner",
			Permission: 0, // for now
		},
		{
			ID:         2,
			Name:       "editor",
			Permission: 0, // for now
		},
		{
			ID:         3,
			Name:       "viewer",
			Permission: 0, // for now
		},
	}

	dbtx.Model(&entities.SpaceRole{})
	dbtx.Create(&spaceRoles)
	if dbtx.Error != nil {
		return dbtx.Error
	}

	return nil
}

func (s *SpaceRoleSeeder) Truncate() error {
	dbtx := databases.GetDB()
	dbtx.Model(&entities.SpaceRole{})
	dbtx.Session(&gorm.Session{AllowGlobalUpdate: true}).Unscoped().Delete(&entities.SpaceRole{})
	if dbtx.Error != nil {
		return dbtx.Error
	}
	return nil
}

package seeders

import (
	"github.com/BlenDMinh/dutgrad-server/databases"
	"github.com/BlenDMinh/dutgrad-server/databases/entities"
)

type SpaceSeeder struct{}

func (s *SpaceSeeder) Name() string {
	return "SpaceSeeder"
}

func (s *SpaceSeeder) Seed() error {
	dbtx := databases.GetDB()
	spaces := []entities.Space{
		{
			Name:          "Test Space 1",
			Description:   "Test Space 1 Description",
			PrivacyStatus: false,
		},
		{
			Name:          "Test Space 2",
			Description:   "Test Space 2 Description",
			PrivacyStatus: true,
		},
		{
			Name:          "Test Space 3",
			Description:   "Test Space 3 Description",
			PrivacyStatus: false,
		},
	}

	dbtx.Model(&entities.Space{}).Create(&spaces)
	if dbtx.Error != nil {
		return dbtx.Error
	}

	return nil
}

func (s *SpaceSeeder) Truncate() error {
	return nil
}

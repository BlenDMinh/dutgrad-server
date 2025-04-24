package seeders

import (
	"github.com/BlenDMinh/dutgrad-server/databases"
	"github.com/BlenDMinh/dutgrad-server/databases/entities"
	"gorm.io/gorm"
)

type TierSeeder struct{}

func (s *TierSeeder) Name() string {
	return "TierSeeder"
}

func (s *TierSeeder) Seed() error {
	dbtx := databases.GetDB()

	tiers := []entities.Tier{
		{
			ID:                1,
			SpaceLimit:        3,
			DocumentLimit:     10,
			QueryHistoryLimit: 50,
			QueryLimit:        30,
			FileSizeLimitKb:   5120,
			ApiCallLimit:      100,
			CostMonth:         0,
			Discount:          0,
		},
		{
			ID:                2,
			SpaceLimit:        10,
			DocumentLimit:     50,
			QueryHistoryLimit: 200,
			QueryLimit:        100,
			FileSizeLimitKb:   15360,
			ApiCallLimit:      500,
			CostMonth:         9.99,
			Discount:          0,
		},
		{
			ID:                3,
			SpaceLimit:        30,
			DocumentLimit:     200,
			QueryHistoryLimit: 1000,
			QueryLimit:        500,
			FileSizeLimitKb:   51200,
			ApiCallLimit:      2000,
			CostMonth:         24.99,
			Discount:          0,
		},
		{
			ID:                4,
			SpaceLimit:        100,
			DocumentLimit:     1000,
			QueryHistoryLimit: 5000,
			QueryLimit:        2000,
			FileSizeLimitKb:   102400,
			ApiCallLimit:      10000,
			CostMonth:         99.99,
			Discount:          0,
		},
	}

	dbtx.Model(&entities.Tier{})
	dbtx.Create(&tiers)
	if dbtx.Error != nil {
		return dbtx.Error
	}

	return nil
}

func (s *TierSeeder) Truncate() error {
	dbtx := databases.GetDB()
	dbtx.Model(&entities.Tier{})
	dbtx.Session(&gorm.Session{AllowGlobalUpdate: true}).Unscoped().Delete(&entities.Tier{})
	if dbtx.Error != nil {
		return dbtx.Error
	}
	return nil
}

package seeders

import (
	"fmt"
	"log"

	"github.com/BlenDMinh/dutgrad-server/databases"
	"github.com/BlenDMinh/dutgrad-server/databases/entities"
	"github.com/BlenDMinh/dutgrad-server/models/dtos"
	"github.com/BlenDMinh/dutgrad-server/services"
)

type MockAccountSeeder struct{}

func (s *MockAccountSeeder) Name() string {
	return "MockAccountSeeder"
}

func (s *MockAccountSeeder) Seed() error {
	authService := new(services.AuthService)
	user, _, _, err := authService.RegisterUser(
		&dtos.RegisterDTO{
			Username: "Test User",
			Email:    "test@example.com",
			Password: "password",
		})
	if err != nil {
		log.Println("Error creating user:", err)
		return err
	}

	dbtx := databases.GetDB()
	tierID := uint(3)
	user.TierID = &tierID
	if err := dbtx.Save(&user).Error; err != nil {
		log.Println("Error assigning tier to user:", err)
		return err
	}

	userSpaces := []entities.Space{}
	for i := 0; i < 3; i++ {
		userSpaces = append(userSpaces, entities.Space{
			Name:          "Own Test Space " + fmt.Sprint(i+1),
			Description:   "Own Test Space " + fmt.Sprint(i+1) + " Description",
			PrivacyStatus: false,
			DocumentLimit: func() int {
				if i == 0 {
					return 20
				} else {
					return 10
				}
			}(),
			FileSizeLimitKb: func() int {
				if i == 0 {
					return 10240
				} else {
					return 5120
				}
			}(),
			ApiCallLimit: func() int {
				if i == 0 {
					return 200
				} else {
					return 100
				}
			}(),
		})
	}

	dbtx.Model(&entities.Space{}).Create(&userSpaces)
	if dbtx.Error != nil {
		log.Println("Error creating spaces:", dbtx.Error)
		return dbtx.Error
	}

	spaceUsers := []entities.SpaceUser{}
	ownerRoleId := uint(1)
	for _, space := range userSpaces {
		spaceUsers = append(spaceUsers, entities.SpaceUser{
			UserID:      user.ID,
			SpaceID:     space.ID,
			SpaceRoleID: &ownerRoleId,
		})
	}

	dbtx.Model(&entities.SpaceUser{}).Create(&spaceUsers)
	if dbtx.Error != nil {
		log.Println("Error creating space users:", dbtx.Error)
		return dbtx.Error
	}

	documents := []entities.Document{}
	for i := 0; i < 3; i++ {
		documents = append(documents, entities.Document{
			SpaceID:       userSpaces[0].ID,
			Name:          "Own Test Document " + fmt.Sprint(i+1),
			PrivacyStatus: false,
			S3URL:         "https://example.com/test-document-" + fmt.Sprint(i+1) + ".pdf",
		})
	}

	dbtx.Model(&entities.Document{}).Create(&documents)
	if dbtx.Error != nil {
		log.Println("Error creating documents:", dbtx.Error)
		return dbtx.Error
	}

	return nil
}

func (s *MockAccountSeeder) Truncate() error {
	return nil
}

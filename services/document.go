package services

import (
	"mime/multipart"

	"github.com/BlenDMinh/dutgrad-server/databases/entities"
	"github.com/BlenDMinh/dutgrad-server/databases/repositories"
)

type DocumentService struct {
	CrudService[entities.Document, uint]
	repo *repositories.DocumentRepository
}

func NewDocumentService() *DocumentService {
	repo := repositories.NewDocumentRepository()
	return &DocumentService{
		CrudService: *NewCrudService(repo),
		repo:        repo,
	}
}

func (s *DocumentService) GetDocumentsBySpaceID(spaceID uint) ([]entities.Document, error) {
	return s.repo.GetBySpaceID(spaceID)
}

func (s *DocumentService) UploadDocument(file *multipart.FileHeader, spaceID uint) (*entities.Document, error) {
	// document := &entities.Document{
	// 	SpaceID: spaceID,
	// 	File:    file,
	// }
	// return s.repo.Create(document)
	return nil, nil
}

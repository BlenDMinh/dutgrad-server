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

func (s *DocumentService) UploadDocument(fileHeader *multipart.FileHeader, spaceID uint) (*entities.Document, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()

	mimeType := fileHeader.Header.Get("Content-Type")
	size := fileHeader.Size

	s3URL, err := UploadFileToS3(fileHeader.Filename, file)
	if err != nil {
		return nil, err
	}

	document := &entities.Document{
		SpaceID:  spaceID,
		Name:     fileHeader.Filename,
		MimeType: mimeType,
		Size:     size,
		S3URL:    s3URL,
	}

	return document, nil
}

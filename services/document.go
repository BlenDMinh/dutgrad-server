package services

import (
	"mime/multipart"

	"github.com/BlenDMinh/dutgrad-server/databases/entities"
	"github.com/BlenDMinh/dutgrad-server/databases/repositories"
	"github.com/BlenDMinh/dutgrad-server/helpers"
)

type DocumentService struct {
	CrudService[entities.Document, uint]
	repo             *repositories.DocumentRepository
	ragServerService *RAGServerService
}

func NewDocumentService() *DocumentService {
	repo := repositories.NewDocumentRepository()
	return &DocumentService{
		CrudService:      *NewCrudService(repo),
		repo:             repo,
		ragServerService: NewRAGServerService(),
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

	// Use the helper to detect proper MIME type
	mimeType, err := helpers.GetMimeType(fileHeader)
	if err != nil {
		return nil, err
	}

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

	document, err = s.repo.Create(document)
	if err != nil {
		return nil, err
	}

	err = s.ragServerService.UploadDocument(fileHeader, spaceID, document.ID)
	if err != nil {
		s.repo.Delete(document.ID)
		return nil, err
	}

	return document, nil
}

func (s *DocumentService) GetUserRoleInSpace(userID, spaceID uint) (string, error) {
	return s.repo.GetUserRoleInSpace(userID, spaceID)
}

func (s *DocumentService) DeleteDocument(documentID uint) error {
	return s.repo.DeleteDocumentByID(documentID)
}

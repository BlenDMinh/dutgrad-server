package services

import (
	"fmt"
	"mime/multipart"

	"github.com/BlenDMinh/dutgrad-server/configs"
	"github.com/BlenDMinh/dutgrad-server/databases"
	"github.com/BlenDMinh/dutgrad-server/databases/entities"
	"github.com/BlenDMinh/dutgrad-server/databases/repositories"
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

func (s *DocumentService) CheckDocumentLimits(spaceID uint, fileSize int64) error {
	db := databases.GetDB()

	var space entities.Space
	if err := db.First(&space, spaceID).Error; err != nil {
		return fmt.Errorf("failed to find space: %v", err)
	}

	var count int64
	if err := db.Model(&entities.Document{}).Where("space_id = ?", spaceID).Count(&count).Error; err != nil {
		return fmt.Errorf("failed to count documents: %v", err)
	}

	if count >= int64(space.DocumentLimit) {
		return fmt.Errorf("document limit reached: this space can only have %d documents", space.DocumentLimit)
	}

	fileSizeKB := fileSize / 1024
	if fileSizeKB > int64(space.FileSizeLimitKb) {
		return fmt.Errorf("file size exceeds the limit of %d KB for this space", space.FileSizeLimitKb)
	}

	return nil
}

func (s *DocumentService) UploadDocument(fileHeader *multipart.FileHeader, spaceID uint, mimeType string) (*entities.Document, error) {
	if err := s.CheckDocumentLimits(spaceID, fileHeader.Size); err != nil {
		return nil, err
	}

	file, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()

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

	env := configs.GetEnv()

	filePath := fmt.Sprintf("%s/documents/view?id=%d", env.WebClientURL, document.ID)

	err = s.ragServerService.UploadDocument(fileHeader, spaceID, document.ID, filePath)
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

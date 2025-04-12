package services

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/BlenDMinh/dutgrad-server/configs"
)

type RAGServerService struct {
	Host              string
	Port              int
	UploadDocumentURL string
}

func NewRAGServerService() *RAGServerService {
	config := configs.GetEnv()
	return &RAGServerService{
		Host:              config.RAGServer.Host,
		Port:              config.RAGServer.Port,
		UploadDocumentURL: config.RAGServer.UploadDocumentURL,
	}
}

func (s *RAGServerService) UploadDocument(fileHeader *multipart.FileHeader, spaceID uint, docId uint) error {
	file, err := fileHeader.Open()
	if err != nil {
		return err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", fileHeader.Filename)
	if err != nil {
		return err
	}

	if _, err = io.Copy(part, file); err != nil {
		return err
	}

	if err = writer.WriteField("spaceId", fmt.Sprintf("%d", spaceID)); err != nil {
		return err
	}
	if err = writer.WriteField("docId", fmt.Sprintf("%d", docId)); err != nil {
		return err
	}

	if err = writer.Close(); err != nil {
		return err
	}

	url := fmt.Sprintf("%s:%d%s", s.Host, s.Port, s.UploadDocumentURL)

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to upload document, status: %d, response: %s", resp.StatusCode, string(respBody))
	}

	return nil
}

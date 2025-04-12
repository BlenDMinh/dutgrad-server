package services

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"

	"github.com/BlenDMinh/dutgrad-server/configs"
	"github.com/BlenDMinh/dutgrad-server/helpers"
)

type RAGServerService struct {
	BaseURL           string
	UploadDocumentURL string
}

func NewRAGServerService() *RAGServerService {
	config := configs.GetEnv()
	return &RAGServerService{
		BaseURL:           config.RAGServer.BaseURL,
		UploadDocumentURL: config.RAGServer.UploadDocumentURL,
	}
}

func (s *RAGServerService) UploadDocument(fileHeader *multipart.FileHeader, spaceID uint, docId uint) error {
	file, err := fileHeader.Open()
	if err != nil {
		return err
	}
	defer file.Close()

	// Get the proper MIME type
	mimeType, err := helpers.GetMimeType(fileHeader)
	if err != nil {
		return err
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Create form file part with proper content type
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="file"; filename="%s"`, fileHeader.Filename))
	h.Set("Content-Type", mimeType)

	part, err := writer.CreatePart(h)
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

	url := fmt.Sprintf("%s%s", s.BaseURL, s.UploadDocumentURL)

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	httpClient := &http.Client{
		Transport: tr,
	}
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

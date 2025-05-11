package services

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
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
	ChatURL           string
	RemoveDocURL      string
	RemoveSpaceURL    string
}

func NewRAGServerService() *RAGServerService {
	config := configs.GetEnv()
	return &RAGServerService{
		BaseURL:           config.RAGServer.BaseURL,
		UploadDocumentURL: config.RAGServer.UploadDocumentURL,
		ChatURL:           config.RAGServer.ChatURL,
		RemoveDocURL:      config.RAGServer.RemoveDocURL,
		RemoveSpaceURL:    config.RAGServer.RemoveSpaceURL,
	}
}

func (s *RAGServerService) UploadDocument(fileHeader *multipart.FileHeader, spaceID uint, docId uint, filePath string) error {
	file, err := fileHeader.Open()
	if err != nil {
		return err
	}
	defer file.Close()

	mimeType, err := helpers.GetMimeType(fileHeader)
	if err != nil {
		return err
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

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
	if err = writer.WriteField("filePath", filePath); err != nil {
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

func (s *RAGServerService) Chat(sessionID uint, spaceID uint, message string) (string, error) {
	url := fmt.Sprintf("%s%s", s.BaseURL, s.ChatURL)

	reqBody := map[string]interface{}{
		"session_id": sessionID,
		"space_id":   spaceID,
		"input":      message,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	httpClient := &http.Client{
		Transport: tr,
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("failed to chat, status: %d, response: %s", resp.StatusCode, string(respBody))
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	var response struct {
		Output string `json:"output"`
	}
	err = json.Unmarshal(respBody, &response)
	if err != nil {
		return "", fmt.Errorf("failed to parse response: %v, raw response: %s", err, string(respBody))
	}

	return response.Output, nil
}

func (s *RAGServerService) RemoveDocument(docId uint, spaceID uint) error {
	url := fmt.Sprintf("%s%s", s.BaseURL, s.RemoveDocURL)

	reqBody := map[string]interface{}{
		"docId":   docId,
		"spaceId": spaceID,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("DELETE", url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

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
		return fmt.Errorf("failed to remove document, status: %d, response: %s", resp.StatusCode, string(respBody))
	}

	return nil
}

func (s *RAGServerService) RemoveSpace(spaceID uint) error {
	url := fmt.Sprintf("%s%s", s.BaseURL, s.RemoveSpaceURL)

	reqBody := map[string]interface{}{
		"spaceId": spaceID,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("DELETE", url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

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
		return fmt.Errorf("failed to remove space, status: %d, response: %s", resp.StatusCode, string(respBody))
	}

	return nil
}

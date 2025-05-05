package helpers

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

func GetUniqueFileKey(filename string) string {
	now := time.Now()

	fileExt := filepath.Ext(filename)

	toHash := fmt.Sprintf("%s-%d-%02d-%02d-%02d-%02d-%02d",
		filename, now.Year(), now.Month(), now.Day(),
		now.Hour(), now.Minute(), now.Second())

	hasher := sha256.New()
	hasher.Write([]byte(toHash))
	hashedKey := hex.EncodeToString(hasher.Sum(nil))

	return hashedKey + fileExt
}

func UploadToS3(bucket string, key string, file multipart.File) (string, error) {
	sess := ConnectAWS()
	s3Client := s3.New(sess)

	uploadParams := &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   file,
	}

	_, err := s3Client.PutObject(uploadParams)
	if err != nil {
		return "", fmt.Errorf("unable to upload file to S3, %v", err)
	}

	fileUrl := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", bucket, key)
	return fileUrl, nil
}

func GetMimeType(fileHeader *multipart.FileHeader) (string, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil && err != io.EOF {
		return "", err
	}

	file.Seek(0, io.SeekStart)

	contentType := http.DetectContentType(buffer)

	if contentType == "application/octet-stream" || contentType == "application/zip" || contentType == "application/x-zip-compressed" {
		ext := filepath.Ext(fileHeader.Filename)
		switch ext {
		case ".pdf":
			return "application/pdf", nil
		case ".doc":
			return "application/msword", nil
		case ".docx":
			return "application/vnd.openxmlformats-officedocument.wordprocessingml.document", nil
		case ".xls":
			return "application/vnd.ms-excel", nil
		case ".xlsx":
			return "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", nil
		case ".ppt":
			return "application/vnd.ms-powerpoint", nil
		case ".pptx":
			return "application/vnd.openxmlformats-officedocument.presentationml.presentation", nil
		case ".txt":
			return "text/plain", nil
		case ".csv":
			return "text/csv", nil
		case ".md":
			return "text/markdown", nil
		}
	}

	return contentType, nil
}

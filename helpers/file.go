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

// GetUniqueFileKey generates a unique key for a file using SHA-256 hash
// It combines the filename and current timestamp, then hashes it
// The file extension is preserved and appended to the hash
func GetUniqueFileKey(filename string) string {
	now := time.Now()

	// Extract file extension
	fileExt := filepath.Ext(filename)

	// Generate a string to hash (filename + timestamp)
	toHash := fmt.Sprintf("%s-%d-%02d-%02d-%02d-%02d-%02d",
		filename, now.Year(), now.Month(), now.Day(),
		now.Hour(), now.Minute(), now.Second())

	// Hash the string using SHA-256
	hasher := sha256.New()
	hasher.Write([]byte(toHash))
	hashedKey := hex.EncodeToString(hasher.Sum(nil))

	// Create the final key with hash and file extension
	return hashedKey + fileExt
}

// UploadToS3 uploads a file to an S3 bucket
func UploadToS3(bucket string, key string, file multipart.File) (string, error) {
	sess := ConnectAWS()
	s3Client := s3.New(sess)

	// Create a new S3 PutObjectInput request
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

// GetMimeType detects the MIME type of a file using content detection
// It takes a file header, opens the file, and detects its MIME type
// Returns the detected MIME type as a string
func GetMimeType(fileHeader *multipart.FileHeader) (string, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Read the first 512 bytes to detect content type
	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil && err != io.EOF {
		return "", err
	}

	// Reset file position
	file.Seek(0, io.SeekStart)

	// Detect content type
	contentType := http.DetectContentType(buffer)

	// For certain file types, http.DetectContentType might not be accurate enough
	// Use file extension as a fallback for common document types
	if contentType == "application/octet-stream" {
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

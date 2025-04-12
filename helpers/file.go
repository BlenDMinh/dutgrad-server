package helpers

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"mime/multipart"
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

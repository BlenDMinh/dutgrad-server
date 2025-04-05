package services

import (
	"fmt"
	"mime/multipart"
	"time"

	"github.com/BlenDMinh/dutgrad-server/helpers"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

func UploadFileToS3(filename string, file multipart.File) (string, error) {
	sess := helpers.ConnectAWS()
	s3Client := s3.New(sess)

	bucket := "dutgrad-doc"

	// Hash from filename and datetime to create a unique key
	now := time.Now()
	key := fmt.Sprintf("%s/%d-%02d-%02d-%02d-%02d-%02d", filename, now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second())

	// Create a new S3 PutObjectInput request
	uploadParams := &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    &key,
		Body:   file,
		ACL:    aws.String("public-read"),
	}

	_, err := s3Client.PutObject(uploadParams)
	if err != nil {
		return "", fmt.Errorf("unable to upload file to S3, %v", err)
	}

	fileUrl := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", bucket, key)
	return fileUrl, nil
}

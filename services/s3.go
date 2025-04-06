package services

import (
	"mime/multipart"

	"github.com/BlenDMinh/dutgrad-server/configs"
	"github.com/BlenDMinh/dutgrad-server/helpers"
)

func UploadFileToS3(filename string, file multipart.File) (string, error) {
	config := configs.GetEnv()

	key := helpers.GetUniqueFileKey(filename)

	return helpers.UploadToS3(config.AWS.S3.Bucket, key, file)
}

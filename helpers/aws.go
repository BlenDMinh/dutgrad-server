package helpers

import (
	"github.com/BlenDMinh/dutgrad-server/configs"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

func ConnectAWS() *session.Session {
	config := configs.GetEnv()
	sess, err := session.NewSession(
		&aws.Config{
			Region: aws.String(config.AWS.Region),
		},
	)
	if err != nil {
		panic(err)
	}
	return sess
}

package helpers

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

func ConnectAWS() *session.Session {
	sess, err := session.NewSession(
		&aws.Config{
			Region:      aws.String("ap-southeast-1"),
			Credentials: credentials.NewEnvCredentials(),
		},
	)
	if err != nil {
		panic(err)
	}
	return sess
}
